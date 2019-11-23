package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/mfslog/sequenceService/infrastructure/config"
	"time"
)

func NewETCDClient(conf config.Config) (cli *clientv3.Client, err error) {
	cfg := clientv3.Config{
		Endpoints:        conf.EtcdAddrs,
		AutoSyncInterval: 1 * time.Second,
		DialTimeout:      5 * time.Second,
	}
	cli, err = clientv3.New(cfg)
	return
}
