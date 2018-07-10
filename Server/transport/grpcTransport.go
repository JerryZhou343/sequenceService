package transport

import (
	"context"
	gcontext "golang.org/x/net/context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	pb "github.com/mfslog/sequenceService/proto"
	"github.com/mfslog/sequenceService/Server/point"
	"github.com/mfslog/sequenceService/Server/service"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/tracing/zipkin"
	stdzipkin "github.com/openzipkin/zipkin-go"
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

func NewGrpcServer(counter metrics.Counter, zipkinTracer *stdzipkin.Tracer) pb.SequenceServer{

	zipkinServer := zipkin.GRPCServerTrace(zipkinTracer)
	options := []grpctransport.ServerOption{
                        zipkinServer,
                    }
	return &grpcServer{
		grpctransport.NewServer(
			sequencePoint.New(service.NewSequenceService(counter),zipkinTracer),
			decodeRequest,
			encodeResponse,
			options...,),
	}
}
