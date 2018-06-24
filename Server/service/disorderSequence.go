package service

import (
	"fmt"
	"time"
)

func GetDisorderSeq(firstId int32, secondId int32, Id int64)string{

	currentTime := time.Now()

	curstr := currentTime.Format("20060102")
	return fmt.Sprintf("%d%d%s%d",firstId,secondId,curstr,Id)
}

