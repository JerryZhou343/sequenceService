package segmentseq_repo

import (
	"github.com/jinzhu/gorm"
	"github.com/mfslog/sequenceService/domain/segment/entity"
	"github.com/mfslog/sequenceService/domain/segment/repository"
	util "github.com/mfslog/sequenceService/infrastructure/repository"
)



type SegmentSeqRepo struct {
	db    *gorm.DB
	table string
}

func NewSegmentSeqRepo(db *gorm.DB) repository.SegmentSeqRepository {
	return &SegmentSeqRepo{
		db:    db,
		table: "t_segment_sequence",
	}
}

func (o *SegmentSeqRepo) Get(pid int32, bid int32) (*entity.SegmentSeq, error) {
	var (
		ret entity.SegmentSeq
		err error
	)
	err =o.db.Table(o.table).Where("product_id = ? and business_id = ? and current_value < max_value").First(&ret).Error
	return &ret, err
}

func (o *SegmentSeqRepo) Save(data *entity.SegmentSeq) (err error) {
	err = o.db.Table(o.table).Update(data).Error
	return err
}

func (o *SegmentSeqRepo) List(page, pageSize int32) ([]*entity.SegmentSeq, int, error) {
	var (
		rets  []*entity.SegmentSeq
		total int
		err error
	)
	o.db.Table(o.table).Count(&total)
	offset, limit := util.CalcPage(page, pageSize)
	err  = o.db.Table(o.table).Offset(offset).Limit(limit).Find(&rets).Error
	return rets, total, err
}

