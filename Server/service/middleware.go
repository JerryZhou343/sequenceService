package service


import (
	gcontext "golang.org/x/net/context"
	pb "github.com/mfslog/sequenceService/proto"
	"github.com/go-kit/kit/metrics"
)

type Middleware func(pb.SequenceServer) pb.SequenceServer

func InstrumentingMiddleware(callCounter metrics.Counter) Middleware{
	return func(next pb.SequenceServer) pb.SequenceServer{
		return instrumentingMiddleware{
			call: callCounter,
			next: next,
		}
	}
}

type instrumentingMiddleware struct{
	call metrics.Counter
	next pb.SequenceServer
}


func (mw instrumentingMiddleware)GetSequence(ctx gcontext.Context, req *pb.SequenceRequest)(*pb.SequenceReply,error){
	v , err := mw.next.GetSequence(ctx, req)
	mw.call.Add(1)
	return v,err
}