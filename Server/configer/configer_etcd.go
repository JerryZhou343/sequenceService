
package configer

import (
    "github.com/coreos/etcd/clientv3"
    "time"
    "context"
    "fmt"
    "strconv"
)

type EtcdService struct {
    
    client *clientv3.Client
}

//连接etcd 服务
//endPoints 为etcd endpoint信息
func (etcd *EtcdService)connectEtcd(endPoints []string)(error){
    var err error
    etcd.client,err = clientv3.New(clientv3.Config{
        Endpoints:   endPoints,
        DialTimeout: 2 * time.Second,
    })
    return err
}


//断开和etcd 服务的连接
func (etcd *EtcdService)disconnect(){
    etcd.client.Close()
}

//通过key值去etcd中查找对应的value
//key 为需要查找的key
//返回了类型为string
func (etcd *EtcdService)getStringValue(key string,defaultValue string)string{
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    
    resp, err := etcd.client.Get(ctx, key)
    cancel()
    if err != nil {
        return ""
    }
    
    var value string = defaultValue
    for _, ev := range resp.Kvs {
        value = fmt.Sprintf("%s",ev.Value)
    }
    
    return value
}


//通过key值去etcd中查找对应的value
//key为需要查找的key
//返回类型为int
func (etcd *EtcdService)getIntVale(key string, defaultValue int)int{
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    
    resp, err := etcd.client.Get(ctx, key)
    cancel()
    if err != nil {
        return 0
    }
    var value int = defaultValue
    for _, ev := range resp.Kvs {
        tmpValue := fmt.Sprintf("%d",ev.Value)
        value,_ = strconv.Atoi(tmpValue)
    }
    
    return value
}