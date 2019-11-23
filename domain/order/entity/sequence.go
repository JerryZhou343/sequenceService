package entity

import (
	"fmt"
	"time"
)

type OrderSeq struct {
	PId           int32 `gorm:"primary_key;column:product_id"`  //复合主键
	BId           int32 `gorm:"primary_key;column:business_id"` // 复合主键
	BaseValue     int64 `gorm:"column:base_value"`
	MaxValue      int64 `gorm:"column:max_value"`
	CurrentValue  int64 `gorm:"column:current_value"`
	StepLength    int32 `gorm:"column:step_length"`
	ResetType     int32 `gorm:"column:reset_type"`
	LastResetTime int32 `gorm:"column:last_reset_time"`
	ProductName  string `gorm:"column:product_name"`
	BusinessName string `gorm:"column:business_name"`
	CreateAt    int64 `gorm:"column:create_at"`
}

func (o *OrderSeq) GetID() string {
	currentTime := time.Now()

	var curResetStamp int32
	//1.获得当前重置时间
	switch o.ResetType { //判断序号reset类型
	case 1: //按天重置
		//1.获得当前日期的零点
		currentTimeStr := currentTime.Format("2006-01-02")
		curDayTime, _ := time.ParseInLocation("2006-01-02", currentTimeStr, time.Local)
		curResetStamp = int32(curDayTime.Unix())
	case 2: //按月重置
		//1.获得xx-01 零点时间
		currentMonthStr := currentTime.Format("2006-01")
		currentMonthTime, _ := time.ParseInLocation("2006-01", currentMonthStr, time.Local)
		curResetStamp = int32(currentMonthTime.Unix())

	case 3: //按年重置
		currentYearStr := currentTime.Format("2006")
		currentYearTime, _ := time.ParseInLocation("2006", currentYearStr, time.Local)
		curResetStamp = int32(currentYearTime.Unix())
	}

	//2.判断重置
	if curResetStamp > o.LastResetTime {
		//重置sequence
		o.CurrentValue = o.BaseValue
		o.LastResetTime = curResetStamp
	}

	newValue := o.CurrentValue + int64(o.StepLength)

	//3. 签发SEQ， 如果超过了最大值限制，则无法签发号
	var sequence string
	if newValue < o.MaxValue {
		curstr := currentTime.Format("20060102")
		o.CurrentValue = o.BaseValue
		sequence = fmt.Sprintf("%d%d%o%.10d", o.PId, o.BId, curstr, newValue)
	}
	return sequence
}
