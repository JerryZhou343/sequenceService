package service

import (
    "context"
    grpc_transport "github.com/go-kit/kit/transport/grpc"
    "github.com/mfslog/sequenceService/proto"
    "github.com/go-kit/kit/endpoint"
    "net"
    "google.golang.org/grpc"
    "fmt"
)


type SequenceServer struct{
    getSequenceHandler grpc_transport.Handler
}


func decodeRequest(_ context.Context, req interface{}) (interface{}, error) {
    return req, nil
}

func encodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
    return rsp, nil
}

func (s *SequenceServer)GetSequence(ctx context.Context, req *sequence.SequenceRequest)(*sequence.SequenceReply,error){
    _, rsp, err := s.getSequenceHandler.ServeGRPC(ctx, req)
    if err != nil {
        return nil, err
    }
    return rsp.(*sequence.SequenceReply),err
}

func makeGetSeqEndpoint() endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        //请求列表时返回 书籍列表
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
    seqServer.getSequenceHandler = seqHandler
    
    ls, _ := net.Listen("tcp", serviceAddress)
    gs := grpc.NewServer(grpc.UnaryInterceptor(grpc_transport.Interceptor))
    
    sequence.RegisterSequenceServer(gs, seqServer)
    
    gs.Serve(ls)
}
