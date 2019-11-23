package rpc

import (
	"github.com/mfslog/sequenceService/infrastructure/config"
	"google.golang.org/grpc"
)

func NewSever(conf *config.Config)*grpc.Server{
	return grpc.NewServer()
}
