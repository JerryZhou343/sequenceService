package application

import (
	"github.com/mfslog/sequenceService/domain/order/entity"
	"github.com/mfslog/sequenceService/domain/snowflake/service"
	"github.com/mfslog/sequenceService/infrastructure/repository/orderseq_repo"
	"github.com/mfslog/sequenceService/infrastructure/repository/segmentseq_repo"
	"github.com/pkg/errors"
	segment "github.com/mfslog/sequenceService/domain/segment/entity"
)

type AppService struct {
	snowSvc service.SnowflakeService
	orderSeqRepo orderseq_repo.OrderRepo
	segmentRepo segmentseq_repo.SegmentSeqRepo


}

func NewAppService(snowSvc service.SnowflakeService) *AppService {
	return &AppService{snowSvc: snowSvc}
}


func (a *AppService)GetOrderID(pid, bid int32)(seq string ,err error){
	var (
		orderSeq *entity.OrderSeq
	)
	orderSeq, err = a.orderSeqRepo.Get(pid,bid)
	if err != nil{
		err = errors.WithStack(err)
		return
	}
	seq = orderSeq.GetID()
	a.orderSeqRepo.Save(orderSeq)
	if err != nil{
		err = errors.WithStack(err)
		return
	}
	return
}

func (a *AppService)GetSnowID()(seq int64, err error){
	seq = a.snowSvc.GetID()
	if seq == 0{
		err = errors.New("invalidate sequence")
	}
	return
}

func (a *AppService)GetDisorderID(pid , bid int32)(seq string, err error){
	seq = a.snowSvc.GetDisorderID(pid, bid)
	if seq == ""{
		err = errors.New("invalidate sequence")
	}
	return
}

func (a *AppService)GetSegmentID(pid, bid int32)(seq int64, err error){
	var (
		segmentSeq *segment.SegmentSeq
	)

	segmentSeq , err = a.segmentRepo.Get(pid, bid)
	if err != nil{
		return 0, err
	}

	seq,err = segmentSeq.GetID()
	return seq , err

}
