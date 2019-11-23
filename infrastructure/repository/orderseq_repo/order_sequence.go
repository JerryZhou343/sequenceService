package orderseq_repo

import (
	"github.com/jinzhu/gorm"
	"github.com/mfslog/sequenceService/domain/order/entity"
	"github.com/mfslog/sequenceService/domain/order/repository"
	util "github.com/mfslog/sequenceService/infrastructure/repository"
)

type OrderRepo struct {
	db    *gorm.DB
	table string
}

func NewOrderRepo(db *gorm.DB) repository.OrderSeqRepository {
	return &OrderRepo{
		db:    db,
		table: "t_order_sequence",
	}
}

func (o *OrderRepo) Get(pid int32, bid int32) (*entity.OrderSeq, error) {
	var (
		ret entity.OrderSeq
		err error
	)
	err = o.db.Table(o.table).Where("product_id = ? and business_id = ?").First(&ret).Error
	return &ret, err
}

func (o *OrderRepo) Save(data *entity.OrderSeq) ( err error) {
	err = o.db.Table(o.table).Update(data).Error
	return
}

func (o *OrderRepo) List(page, pageSize int32) ([]*entity.OrderSeq, int, error) {
	var (
		rets  []*entity.OrderSeq
		total int
		err error
	)
	o.db.Table(o.table).Count(&total)
	offset, limit := util.CalcPage(page, pageSize)
	err = o.db.Table(o.table).Offset(offset).Limit(limit).Find(&rets).Error
	return rets, total, err
}
