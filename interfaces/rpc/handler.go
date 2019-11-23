package rpc

import (
	"context"
	"github.com/mfslog/sequenceService/application"
	sequence "github.com/mfslog/sequenceService/proto"
)

type  Handler struct{
	app * application.AppService
}

func Newhandler(app *application.AppService)*Handler{
	return &Handler{app:app}
}


func (h *Handler)GetOrderID(ctx context.Context, req *sequence.GetOrderIDReq) (rsp *sequence.GetOrderIDRsp,err error){
	var (
		id string
	)
	rsp = &sequence.GetOrderIDRsp{}
	id , err = h.app.GetOrderID(req.Pid, req.Bid)
	rsp.Id = id
	return
}

func (h *Handler)GetSnowflakeID(ctx context.Context, req *sequence.GetSnowflakeIDReq)(rsp *sequence.GetSnowflakeIDRsp,err error){
	var (
		id int64
	)
	rsp = &sequence.GetSnowflakeIDRsp{}

	id , err = h.app.GetSnowID()
	rsp.Id = id
	return

}

func (h *Handler)GetDisorderID(ctx context.Context, req *sequence.GetDisorderIDReq)(rsp *sequence.GetDisorderIDRsp,err error){
	var (
		id string
	)

	rsp = &sequence.GetDisorderIDRsp{}

	id, err  = h.app.GetDisorderID(req.Pid, req.Bid)
	rsp.Id = id
	return
}

func (h *Handler)GetSegmentID(ctx context.Context, req *sequence.GetSegmentIDReq)( rsp *sequence.GetSegmentIDRsp,err error){
	var (
		id int64
	)

	id, err = h.app.GetSegmentID(req.Pid, req.Bid)
	rsp.Id = id
	return
}
