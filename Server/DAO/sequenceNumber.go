package DAO

import (
    "github.com/mfslog/sequenceService/Server/DBSession"
    "github.com/mfslog/sequenceService/Server/cacheSession"
    "fmt"
    "github.com/mfslog/sequenceService/Server/log"
    "strconv"
)

type TSequenceNumber struct{
    FirstId int32 `gorm:"primary_key;column:first_id"`
    SecondId int32 `gorm:"primary_key;column:second_id"`
    BaseValue int64 `gorm:"column:base_value"`
    MaxValue int64  `gorm:"column:max_value"`
    CurrentValue int64 `gorm:"column:current_value"`
    StepLength int32 `gorm:"column:step_length"`
    ResetType int32 `gorm:"column:reset_type"`
    LastResetTime int32 `gorm:"column:last_reset_time"`
}

//指定表名
func (TSequenceNumber)TableName()string{
    return "t_sequence_number"
}


//查询序列号表
func (seq *TSequenceNumber)GetOneByBusinessID(firstId int32, secondId int32){
    seq.getOneByBussinessIDFromCache(firstId,secondId)
}

//重置序列号
func (seq *TSequenceNumber)ResetSeqByBussinessID(firstId int32, secondId int32,lastResetTime int32){
    seq.resetSeqByBussinessIDFromCache(firstId,secondId,lastResetTime)
    seq.resetSeqByBussinessIDFromDB(firstId,secondId,lastResetTime)
}

//更新序列号
func (seq *TSequenceNumber)UpdateSeqByBussinessID(firstId int32, secondId int32,newValue int64){
    seq.updateSeqByBussinessIDFromCache(firstId,secondId,newValue)
    seq.updateSeqByBusinessIDFromDB(firstId,secondId,newValue)
}

//查询序列号from mysql db
func (seq *TSequenceNumber)getOneByBusinessIDFromDB(firstId int32 , secondId int32){
    // 获得数据库链接实例
    db := DBSession.GetInstance()
    // 查询
    
    db.SYSDB.Select("first_id,second_id,base_value,max_value,current_value,step_length,reset_type,last_reset_time").
        Where("first_id = ? and second_id = ? ",firstId,secondId).First(seq)
}


//重置序列号 from mysql db
func (seq *TSequenceNumber)resetSeqByBussinessIDFromDB(firstId int32, secondId int32,lastResetTime int32){
    db := DBSession.GetInstance()
    db.SYSDB.Model(seq).Where("first_id = ? and second_id = ?",firstId,secondId).
        Updates(map[string]interface{}{"current_value":seq.BaseValue, "last_reset_time":lastResetTime})
}


//更新序号值 from mysql db
func (seq *TSequenceNumber)updateSeqByBusinessIDFromDB(firstId int32, secondId int32, currentValue int64){
    db := DBSession.GetInstance()
    db.SYSDB.Model(seq).Where("first_id = ? and second_id = ?",firstId,secondId).
        Update("current_value",currentValue)
}


//查询序列号from redis db
func (seq *TSequenceNumber)getOneByBussinessIDFromCache(firstId int32,secondId int32){

    //捕获查询异常
    defer func(firstId int32, secondId int32){
        if err := recover(); err != nil{
            seq.getOneByBusinessIDFromDB(firstId,secondId)
        }
    }(firstId,secondId)

    ins := cacheSession.GetInstance()
    client := ins.GetClient()
    key := fmt.Sprintf("t_seq_%d_%d",firstId,secondId)
    keyExists, _ := client.Exists(key).Result()

    if keyExists > 0{
        values, err := client.HMGet(key,"current_value","step_length","last_reset_time",
            "reset_type","base_value","max_value").Result()
        if err == nil && values != nil {
            var tmp int
            seq.FirstId = firstId
            seq.SecondId = secondId
            seq.CurrentValue,_ = strconv.ParseInt(values[0].(string),10,64)

            tmp,_ = strconv.Atoi(values[1].(string))
            seq.StepLength = int32(tmp)

            tmp,_ = strconv.Atoi(values[2].(string))
            seq.LastResetTime =int32(tmp)

            tmp, _ = strconv.Atoi(values[3].(string))
            seq.ResetType = int32(tmp)

            seq.BaseValue ,_= strconv.ParseInt(values[4].(string),10,64)

            seq.MaxValue,_ = strconv.ParseInt(values[5].(string),10,64)
            //fmt.Println(seq.CurrentValue)
        }
    } else{
        seq.getOneByBusinessIDFromDB(firstId,secondId)
    }
}

//重置序列号 from redis db
func (seq *TSequenceNumber)resetSeqByBussinessIDFromCache(firstId int32, secondId int32,lastResetTime int32){
    ins := cacheSession.GetInstance()
    client := ins.GetClient()
    key := fmt.Sprintf("t_seq_%d_%d",firstId,secondId)
    _, err := client.
        HMSet(key,map[string]interface{}{"current_value":seq.BaseValue,"step_length":seq.StepLength,
        "last_reset_time":lastResetTime,"reset_type":seq.ResetType,"base_value":seq.BaseValue,
        "max_value":seq.MaxValue}).
        Result()

    if err != nil{
        log.Error("update value to redis failed.")
    }
}

//更新序号值 from redis db
func (seq *TSequenceNumber)updateSeqByBussinessIDFromCache(firstId int32,secondId int32, newValue int64){
    ins := cacheSession.GetInstance()
    client := ins.GetClient()
    key := fmt.Sprintf("t_seq_%d_%d",firstId,secondId)
    _, err := client.
        HMSet(key,map[string]interface{}{"current_value":newValue,"step_length":seq.StepLength,
        "last_reset_time":seq.LastResetTime,"reset_type":seq.ResetType,"base_value":seq.BaseValue,
        "max_value":seq.MaxValue}).
        Result()

    if err != nil{
        log.Error("update value to redis failed.")
    }
}

//add new seq to redis
func (seq *TSequenceNumber)addSeqToCache(firstId int32, secondId int32){
    ins := cacheSession.GetInstance()
    client := ins.GetClient()
    key := fmt.Sprintf("t_seq_%d_%d",firstId,secondId)
    _, err := client.
        HMSet(key,map[string]interface{}{"current_value":seq.CurrentValue,"step_length":seq.StepLength,
        "last_reset_time":seq.LastResetTime,"reset_type":seq.ResetType,"base_value":seq.BaseValue,
        "max_value":seq.MaxValue}).
        Result()

    if err != nil{
        log.Error("update value to redis failed.")
    }
}
