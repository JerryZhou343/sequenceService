package serverplugin

import (
    "github.com/go-kit/kit/sd/etcdv3"
    "time"
    "context"
    "strings"
    "strconv"
    "os"
    "github.com/go-kit/kit/log"
    "github.com/mfslog/sequenceService/Server/common"
)

type ServiceInfo struct{
    Name string
    BasePath string
    IP string
    Port int
}


type EtcdService struct{
    Info ServiceInfo
    Client etcdv3.Client
}



func NewEtcdService(info ServiceInfo, endPoints []string)(*EtcdService, error){
    options := etcdv3.ClientOptions{
        DialTimeout: time.Second * 3,
        DialKeepAlive: time.Second * 3,
    }
    ctx        := context.Background()
    //创建etcd连接
    cli, err := etcdv3.NewClient(ctx, endPoints, options)
    if err != nil {
        os.Exit(1)
    }
    
    return &EtcdService{
        Info:info,
        Client:cli,
    },err
}


func (s *EtcdService)RegisterService(){
    prefixKey := strings.TrimRight(s.Info.BasePath,"/") + "/" + s.Info.Name +"/" + "host_list_t/"
    fPrefixKey := strings.TrimRight(s.Info.BasePath,"/") + "/" + s.Info.Name +"/" + "host_list_f/"
    instance := s.Info.IP +":" + strconv.Itoa(s.Info.Port)
    key := prefixKey + instance

    //记录常驻机器地址
    common.HostFPath = fPrefixKey + instance
    value := instance
    // 创建注册器
    registrar := etcdv3.NewRegistrar(s.Client, etcdv3.Service{
        Key:   key,
        Value: value,
        TTL :etcdv3.NewTTLOption(3,5),
    }, log.NewNopLogger())
    // 注册器启动注册
    registrar.Register()
}