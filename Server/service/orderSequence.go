package service

import (
    "github.com/mfslog/sequenceService/Server/DAO"
    "time"
    "fmt"
    "github.com/mfslog/sequenceService/Server/log"
)

func GetOrderSequnce(firstId int32, secondId int32)string{
    seqRecord := DAO.TSequenceNumber{}
    seqRecord.GetOneByBusinessID(firstId, secondId)
    currentTime := time.Now()
    //判断序号reset类型
    var lastRestTimeStamp int64 = int64(seqRecord.LastResetTime)
    lastResetTime := time.Unix( lastRestTimeStamp,0)
    switch(seqRecord.RsetType){
    case 1: //按天重置
        if currentTime.After(lastResetTime.Add(time.Hour*24)) {
            currentTimeStr := currentTime.Format("yyyy-mm-dd")
            curDayTime,_ := time.Parse("yyyy-mm-dd",currentTimeStr)
            curDayStamp := curDayTime.Unix()
            var tmp int32 = int32(curDayStamp)
            //重置sequence
            seqRecord.ResetSeqByBussinessID(firstId,secondId,tmp)
            //重新获得sequence值
            seqRecord.GetOneByBusinessID(firstId,secondId)
        }
    }
    
    newValue := seqRecord.CurrentValue + int64(seqRecord.StepLength)
    var sequence string
    
    if newValue < seqRecord.MaxValue {
        curstr := currentTime.Format("yyyymmdd")
        seqRecord.UpdateSeqByBusinessId(firstId,secondId,newValue)
        sequence = fmt.Sprintf("%s%d%d%d",curstr,firstId,secondId,newValue)
    }else{
        log.Error(0,fmt.Sprintf("error current sequence for [%d:%d] Depletion",firstId,secondId))
        sequence = ""
    }
    
    return sequence
}