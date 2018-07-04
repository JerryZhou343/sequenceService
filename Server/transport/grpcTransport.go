package transport

import (
	"context"
	gcontext "golang.org/x/net/context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	pb "github.com/mfslog/sequenceService/proto"
	"github.com/mfslog/sequenceService/Server/point"
	"github.com/mfslog/sequenceService/Server/service"
	"fmt"
	"github.com/mfslog/sequenceService/Server/log"
	"net"
	"google.golang.org/grpc"
)



type grpcServer struct{
	getSequence grpctransport.Handler
}

func decodeRequest(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func encodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
	return rsp, nil
}


func (s* grpcServer)GetSequence(ctx gcontext.Context, req *pb.SequenceRequest)(*pb.SequenceReply,error){
	_, rep, err := s.getSequence.ServeGRPC(ctx,req)
	if err != nil{
		return nil, err
	}

	return rep.(*pb.SequenceReply),nil
}

func NewGrpcServer() pb.SequenceServer{
	return &grpcServer{
		grpctransport.NewServer(
			sequencePoint.MakeGetSequenceEndpoint(service.NewSequenceService()),
			decodeRequest,
			encodeResponse),
	}
}

func NewServer(port int){
	serviceAddress := fmt.Sprintf("0.0.0.0:%d",port)
	log.Info("listen:" + serviceAddress)

	ls, _ := net.Listen("tcp", serviceAddress)
	gs := grpc.NewServer()

	pb.RegisterSequenceServer(gs, NewGrpcServer())

	gs.Serve(ls)
}