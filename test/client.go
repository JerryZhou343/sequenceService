package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"github.com/mfslog/sequenceService/proto"
	"google.golang.org/grpc"
	"io"
	"os"
	"time"
)

var logger log.Logger

func NewLogger() log.Logger {
	var fd *os.File
	fd, _ = os.Create("seq.log")
	logger := log.NewJSONLogger(log.NewSyncWriter(fd))
	logger = log.With(logger,
		//"module", MineDodule,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)
	logger = level.NewFilter(logger, level.AllowAll())
	return logger
}

func main() {
	logger = NewLogger()
	var (
		//注册中心地址
		etcdServer = "192.168.0.109:2379"
		//监听的服务前缀
		prefix = "/private.com/service_cluster/sequenceNumber/host_list_t"
		ctx    = context.Background()
	)
	options := etcdv3.ClientOptions{
		DialTimeout:   time.Second * 3,
		DialKeepAlive: time.Second * 3,
	}
	//连接注册中心
	client, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil {
		panic(err)
	}
	//logger := log.NewNopLogger()
	//创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
	instancer, err := etcdv3.NewInstancer(client, prefix, logger)
	if err != nil {
		panic(err)
	}
	//创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
	endpointer := sd.NewEndpointer(instancer, reqFactory, logger)
	//创建负载均衡器
	balancer := lb.NewRoundRobin(endpointer)

	/**
	  我们可以通过负载均衡器直接获取请求的endPoint，发起请求
	  reqEndPoint,_ := balancer.Endpoint()
	*/

	/**
	  也可以通过retry定义尝试次数进行请求
	*/
	reqEndPoint := lb.Retry(3, 3*time.Second, balancer)

	//现在我们可以通过 endPoint 发起请求了
	req := struct{}{}
	for {
		//fmt.Println(time.Now())
		logger.Log("msg", time.Now().String())
		if _, err = reqEndPoint(ctx, req); err != nil {
			//panic(err)
			//fmt.Println(err)
			logger.Log("msg", err.Error())
		}

		logger.Log("msg", time.Now().String())
		//fmt.Println(time.Now())
	}

}

//通过传入的 实例地址  创建对应的请求endPoint
func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//fmt.Println("请求服务: ", instanceAddr)
		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			fmt.Println(err)
			panic("connect error")
		}
		defer conn.Close()
		seqClient := sequence.NewSequenceClient(conn)
		stamp := time.Now().UnixNano()
		seqInfo, _ := seqClient.GetSequence(context.Background(), &sequence.SequenceRequest{CallSeq: stamp, FirstBID: 42000000, SecondBID: 58, Target: 3, Mode: 1})
		fmt.Println("==========get==========")
		fmt.Println("origin stamp", stamp)
		fmt.Println("call seq", seqInfo.GetCallSeq())
		fmt.Println("call id:", seqInfo.GetCallID())
		fmt.Println("order:", seqInfo.GetSeq())
		fmt.Println("==========end==========")
		logger.Log("msg", "==================")
		logger.Log("msg", seqInfo.GetCallID())
		logger.Log("msg", seqInfo.GetSeq())
		logger.Log("msg", "==================")
		return nil, nil
	}, nil, nil
}
