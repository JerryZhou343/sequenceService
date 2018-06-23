package DAO

import (
    "github.com/mfslog/sequenceService/Server/DBSession"
)

type TSequenceNumber struct{
    FirstId int32 `gorm:"primary_key;column:first_id"`
    SecondId int32 `gorm:"primary_key;column:second_id"`
    BaseValue int64 `gorm:"column:base_value"`
    MaxValue int64  `gorm:"column:max_value"`
    CurrentValue int64 `gorm:"column:current_value"`
    StepLength int32 `gorm:"column:step_length"`
    RsetType int32 `gorm:"column:reset_type"`
    LastResetTime int32 `gorm:"column:last_reset_time"`
}

//指定表名
func (TSequenceNumber)TableName()string{
    return "t_sequence_number"
}


//查询序列号表
func (seq *TSequenceNumber)GetOneByBusinessID(firstId int32 , secondId int32){
    // 获得数据库链接实例
    db := DBSession.GetInstance()
    // 查询
    
    db.SYSDB.Select("first_id,second_id,base_value,max_value,current_value,step_length,reset_type,last_reset_time").
        Where("first_id = ? and second_id = ? ",firstId,secondId).First(seq)
}


//重置序列号
func (seq *TSequenceNumber)ResetSeqByBussinessID(firstId int32, secondId int32,lastResetTime int32){
    db := DBSession.GetInstance()
    db.SYSDB.Model(seq).Where("first_id = ? and second_id = ?",firstId,secondId).
        Updates(map[string]interface{}{"current_value":seq.BaseValue, "last_reset_time":lastResetTime})
}


//更新序号值
func (seq *TSequenceNumber)UpdateSeqByBusinessId(firstId int32, secondId int32, currentValue int64){
    db := DBSession.GetInstance()
    db.SYSDB.Model(seq).Where("first_id = ? and second_id = ?",firstId,secondId).
        Update("current_value",currentValue)
}