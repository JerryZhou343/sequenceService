package entity

import "github.com/mfslog/sequenceService/domain/errors"

type SegmentSeq struct {
	ID           int64 `gorm:"AUTO_INCREMENT;primary_key"`
	PId          int32 `gorm:"index;column:product_id"`  //复合主键
	BId          int32 `gorm:"index;column:business_id"` // 复合主键
	BaseValue    int64 `gorm:"column:base_value"`
	MaxValue     int64 `gorm:"column:max_value"`
	CurrentValue int64 `gorm:"column:current_value"`
	StepLength   int32 `gorm:"column:step_length"`
	ProductName  string `gorm:"column:product_name"`
	BusinessName string `gorm:"column:business_name"`
	CreateAt    	int64 `gorm:"column:create_at"`
}

func (s *SegmentSeq) GetID() (int64,error) {
	newValue := s.CurrentValue + int64(s.StepLength)
	s.CurrentValue = newValue
	if newValue > s.MaxValue{
		return 0,errors.ErrSeqNotEnough
	}
	return s.CurrentValue, nil
}
