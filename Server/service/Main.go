package service

import (
    "context"
    grpc_transport "github.com/go-kit/kit/transport/grpc"
    gcontext "golang.org/x/net/context"
    pb "github.com/mfslog/sequenceService/proto"
    "github.com/go-kit/kit/endpoint"
    "net"
    "google.golang.org/grpc"
    "fmt"
    "github.com/mfslog/sequenceService/Server/log"
    "github.com/go-kit/kit/metrics"
)

var (
    snowFlake SnowFlake
)


type basicService struct{}

//服务
func (s basicService)GetSequence(ctx gcontext.Context, req *pb.SequenceRequest)(*pb.SequenceReply,error){
        seq := new(pb.SequenceReply)
        seq.CallSeq = req.CallSeq

        id := snowFlake.GetSnowflakeId()
        if req.Target == 1 || req.Target == 3{
            seq.CallID = id
        }

        if req.Target == 2 || req.Target == 3{
            if req.Mode == 1{
                // 有序序号
                seq.Seq = GetOrderSequence(req.FirstBID,req.SecondBID)
            }else{
                // 无序序号
                seq.Seq = GetDisorderSeq(req.FirstBID,req.SecondBID,id)
            }
        }
        return seq,nil
}

func NewBasicService() pb.SequenceServer{
    return basicService{}
}

type ServiceMddileware func(server pb.SequenceServer)pb.SequenceServer

func NewSequenceService( counter metrics.Counter, chart metrics.Histogram)pb.SequenceServer{
    var svc pb.SequenceServer
    {
        svc = NewBasicService()
        svc = instrumentingMiddleware(counter,chart)(svc)

    }
    return svc
}







func decodeRequest(_ context.Context, req interface{}) (interface{}, error) {
    return req, nil
}



func encodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
    return rsp, nil
}

//建立Seq 服务接入点
func makeGetSeqEndpoint(s SequenceService) endpoint.Endpoint {
    return s.GetSequence()
}



func NewServer(port int){
    
    //构建服务
    seqServer := new(SequenceService)
    //创建一个传输层
    seqHandler := grpc_transport.NewServer(
        makeGetSeqEndpoint(),
        decodeRequest,
        encodeResponse,
    )
    snowFlake.Init()
    //监听服务
    serviceAddress := fmt.Sprintf("0.0.0.0:%d",port)
    log.Info("listen:" + serviceAddress)
    seqServer.getSequenceHandler = seqHandler
    
    ls, _ := net.Listen("tcp", serviceAddress)
    gs := grpc.NewServer()
    
    sequence.RegisterSequenceServer(gs, seqServer)
    
    gs.Serve(ls)
}
