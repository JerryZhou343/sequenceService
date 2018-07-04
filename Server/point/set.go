package sequencePoint

import ("github.com/go-kit/kit/endpoint"
	pb "github.com/mfslog/sequenceService/proto"
	gcontext "golang.org/x/net/context"
)



func MakeGetSequenceEndpoint(s pb.SequenceServer) endpoint.Endpoint{
	return func(ctx gcontext.Context, request interface{})(response interface{},err error){
		req,_ := request.(pb.SequenceRequest)
		return s.GetSequence(ctx,&req)
	}
}