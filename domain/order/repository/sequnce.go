package repository

import "github.com/mfslog/sequenceService/domain/order/entity"

type OrderSeqRepository interface {
	Get(pid, bid int32) (*entity.OrderSeq, error)
	Save(data *entity.OrderSeq) error
	List(page, pageSize int32) ([]*entity.OrderSeq, int, error)
}
