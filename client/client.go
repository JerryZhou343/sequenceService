package main

import (
    "context"
    "github.com/go-kit/kit/sd/etcdv3"
    "time"
    "github.com/go-kit/kit/sd"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/endpoint"
    "io"
    "github.com/go-kit/kit/sd/lb"
    "fmt"
    "google.golang.org/grpc"
    "github.com/mfslog/sequenceService/proto"
)

func main() {
    
    var (
        //注册中心地址
        etcdServer = "127.0.0.1:2379"
        //监听的服务前缀
        prefix     = "/private.com/service_cluster/sequenceNumber/host_list_t"
        ctx        = context.Background()
    )
    options := etcdv3.ClientOptions{
        DialTimeout: time.Second * 3,
        DialKeepAlive: time.Second * 3,
    }
    //连接注册中心
    client, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
    if err != nil {
        panic(err)
    }
    logger := log.NewNopLogger()
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
    if _, err = reqEndPoint(ctx, req); err != nil {
        panic(err)
    }
}


//通过传入的 实例地址  创建对应的请求endPoint
func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        fmt.Println("请求服务: ", instanceAddr)
        conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
        if err != nil {
            fmt.Println(err)
            panic("connect error")
        }
        defer conn.Close()
        seqClient := sequence.NewSequenceClient(conn)
        seqInfo,_:=seqClient.GetSequence(context.Background(),&sequence.SequenceRequest{CallSeq: 1, FirstBID:420000000, SecondBID:58,Target:3, Mode:2})
        fmt.Println("get call id")
        fmt.Println("call seq", " => ", "call id:", seqInfo.GetCallID())
        fmt.Println("sequence " , seqInfo.GetCallSeq() )
        return nil,nil
    },nil,nil
}
