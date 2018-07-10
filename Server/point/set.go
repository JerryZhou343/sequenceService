package sequencePoint

import ("github.com/go-kit/kit/endpoint"
	pb "github.com/mfslog/sequenceService/proto"
	gcontext "golang.org/x/net/context"
	stdZipkin "github.com/openzipkin/zipkin-go"
	"github.com/go-kit/kit/tracing/zipkin"
)



func MakeGetSequenceEndpoint(s pb.SequenceServer) endpoint.Endpoint{
	return func(ctx gcontext.Context, request interface{})(response interface{},err error){
		req,_ := request.(*pb.SequenceRequest)
		return s.GetSequence(ctx,req)
	}
}

func New(service  pb.SequenceServer, tracer *stdZipkin.Tracer) endpoint.Endpoint{
	point := MakeGetSequenceEndpoint(service)
	return zipkin.TraceEndpoint(tracer,"sequence")(point)
}