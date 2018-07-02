package service

import (
	"github.com/go-kit/kit/metrics"
	"time"
	"fmt"
	pb "github.com/mfslog/sequenceService/proto"
	gcontext "golang.org/x/net/context"
)

func instrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
)ServiceMddileware{
	return func(next pb.SequenceServer)pb.SequenceServer{
		return instrmw{requestCount, requestLatency, next}
	}
}

type instrmw struct{
	requestCount metrics.Counter
	requestLatency metrics.Histogram
	next pb.SequenceServer
}


func(mw instrmw)GetSequence(ctx gcontext.Context, req *pb.SequenceRequest)(output *pb.SequenceReply,err error){
	defer func(begin time.Time) {
		lvs := []string{"method","getSequence","error",fmt.Sprint( err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output,err = mw.next.GetSequence(ctx, req)
	return
}
