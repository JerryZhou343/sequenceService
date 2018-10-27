package service

import (
    gcontext "golang.org/x/net/context"
    pb "github.com/mfslog/sequenceService/proto"
    "github.com/go-kit/kit/metrics"
    "github.com/mfslog/sequenceService/Server/log"
    "fmt"
)



type basicService struct{
    snowFlake SnowFlake
}

//服务
func (s basicService)GetSequence(ctx gcontext.Context, req *pb.SequenceRequest)(*pb.SequenceReply,error){
    seq := new(pb.SequenceReply)
    seq.CallSeq = req.CallSeq

    id := s.snowFlake.GetSnowflakeId()
    log.Info(fmt.Sprintf("accept call request: %s",seq.CallSeq))
    seq.Seq = id
    return seq,nil
}

func NewBasicService() pb.SequenceServer{
    var snowFlake SnowFlake
    snowFlake.Init()
    return basicService{snowFlake}
}


func NewSequenceService(callCounter metrics.Counter )pb.SequenceServer{
    var svc pb.SequenceServer
    {
        svc = NewBasicService()
        svc = InstrumentingMiddleware(callCounter)(svc)
    }
    return svc
}

