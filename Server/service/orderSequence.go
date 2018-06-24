package service

import (
    "github.com/mfslog/sequenceService/Server/DAO"
    "time"
    "fmt"
    "github.com/mfslog/sequenceService/Server/log"
    "github.com/mfslog/sequenceService/Server/serverPlugin"
    "github.com/mfslog/sequenceService/Server/common"
)

func GetOrderSequence(firstId int32, secondId int32)string{

    etcd := serverplugin.GetEtcdConIns()
    r := etcd.Lock(common.LockPath)
    if r == common.R_ERR{
        return ""
    }

    defer etcd.UnLock()

    seqRecord := DAO.TSequenceNumber{}
    seqRecord.GetOneByBusinessID(firstId, secondId)
    currentTime := time.Now()

    var curResetStamp int32
    //1.获得当前重置时间
    switch(seqRecord.ResetType){//判断序号reset类型
    case 1: //按天重置
        //1.获得当前日期的零点
        currentTimeStr := currentTime.Format("2006-01-02")
        curDayTime,_ := time.ParseInLocation("2006-01-02",currentTimeStr,time.Local)
        curResetStamp = int32(curDayTime.Unix())
    case 2: //按月重置
        //1.获得xx-01 零点时间
        currentMonthStr := currentTime.Format("2006-01")
        currentMonthTime,_ := time.ParseInLocation("2006-01",currentMonthStr,time.Local)
        curResetStamp = int32(currentMonthTime.Unix())

    case 3: //按年重置
        currentYearStr := currentTime.Format("2006")
        currentYearTime,_ := time.ParseInLocation("2006",currentYearStr,time.Local)
        curResetStamp = int32(currentYearTime.Unix())
    }

    //2.判断重置
    if curResetStamp> seqRecord.LastResetTime {
        //重置sequence
        seqRecord.ResetSeqByBussinessID(firstId,secondId,curResetStamp)
        //重新获得sequence值
        seqRecord.GetOneByBusinessID(firstId,secondId)
    }

    newValue := seqRecord.CurrentValue + int64(seqRecord.StepLength)

    var sequence string
    if newValue < seqRecord.MaxValue {
        curstr := currentTime.Format("20060102")
        seqRecord.UpdateSeqByBussinessID(firstId,secondId,newValue)
        sequence = fmt.Sprintf("%d%d%s%.10d",firstId,secondId,curstr,newValue)
    }else{
        log.Error(fmt.Sprintf("error current sequence for [%d:%d] value [%d:%d] Depletion",
            firstId,secondId,newValue, seqRecord.MaxValue))
        sequence = ""
    }
    return sequence
}