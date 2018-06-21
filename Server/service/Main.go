package service

import (
    "context"
    grpc_transport "github.com/go-kit/kit/transport/grpc"
    gcontext "golang.org/x/net/context"
    "github.com/mfslog/sequenceService/proto"
    "github.com/go-kit/kit/endpoint"
    "net"
    "google.golang.org/grpc"
    "fmt"
    "github.com/mfslog/sequenceService/Server/log"
)


type SequenceServer struct{
    getSequenceHandler grpc_transport.Handler
}


func (s *SequenceServer)GetSequence(ctx gcontext.Context, req *sequence.SequenceRequest)(*sequence.SequenceReply,error){
    _, rsp, err := s.getSequenceHandler.ServeGRPC(ctx, req)
    if err != nil {
        return nil, err
    }
    return rsp.(*sequence.SequenceReply),err
}


func decodeRequest(_ context.Context, req interface{}) (interface{}, error) {
    return req, nil
}



func encodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
    return rsp, nil
}


func makeGetSeqEndpoint() endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(*sequence.SequenceRequest)
        seq := new(sequence.SequenceReply)
        seq.CallSeq = req.CallSeq
        seq.CallID = GetSnowflakeId()
        return seq,nil
    }
}



func NewServer(port int){
    
    //构建服务
    seqServer := new(SequenceServer)
    seqHandler := grpc_transport.NewServer(
        makeGetSeqEndpoint(),
        decodeRequest,
        encodeResponse,
    )
    
    //监听服务
    serviceAddress := fmt.Sprintf("0.0.0.0:%d",port)
    log.Info(0,"listen:" + serviceAddress)
    seqServer.getSequenceHandler = seqHandler
    
    ls, _ := net.Listen("tcp", serviceAddress)
    gs := grpc.NewServer()
    
    sequence.RegisterSequenceServer(gs, seqServer)
    
    gs.Serve(ls)
}
