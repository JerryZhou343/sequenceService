package service

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/mfslog/sequenceService/infrastructure/config"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type SnowflakeService interface {
	GetID() int64
	GetDisorderID(pid, bid int32) string
}

type snowflakeService struct {
	key           string // etcd key
	machineID     int64  // 机器 id 占10位, 十进制范围是 [ 0, 1023 ]
	sn            int64  // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
	lastTimeStamp int64  // 上次的时间戳(毫秒级), 1秒=1000毫秒, 1毫秒=1000微秒,1微秒=1000纳秒
	etcdCli       *clientv3.Client
}

func NewSnowflakeService(conf *config.Config, etcdCli *clientv3.Client) (SnowflakeService, error) {
	svc := &snowflakeService{
		key: conf.EtcdPath + "/" + "snowflakeService",
		// 把机器 id 左移 12 位,让出 12 位空间给序列号使用
		machineID: int64(conf.SrvID) << 12,
		etcdCli:   etcdCli,
	}
	err := svc.getLastTimeStamp()

	return svc, err
}

func (s *snowflakeService) getLastTimeStamp() error {
	var (
		lastTimeStamp int64
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	rsp, err := s.etcdCli.Get(ctx, s.key)
	defer cancel()
	if err != nil {
		return errors.WithStack(err)
	}
	for _, ev := range rsp.Kvs {
		tmpValue := fmt.Sprintf("%d", ev.Value)
		lastTimeStamp, _ = strconv.ParseInt(tmpValue, 10, 64)
	}

	if lastTimeStamp == 0 {
		lastTimeStamp = time.Now().UnixNano() / 1000000
	}
	s.lastTimeStamp = lastTimeStamp
	return nil
}

func (s *snowflakeService) setLastTimeStamp() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = s.etcdCli.Put(ctx, s.key, fmt.Sprintf("%d", s.lastTimeStamp))
	defer cancel()
	if err != nil {
		return err
	}

	return
}

func (s *snowflakeService) GetID() int64 {
	curTimeStamp := time.Now().UnixNano() / 1000000

	// 同一毫秒
	if curTimeStamp == s.lastTimeStamp {
		s.sn++
		// 序列号占 12 位,十进制范围是 [ 0, 4095 ]
		if s.sn > 4095 {
			time.Sleep(time.Millisecond)
			curTimeStamp = time.Now().UnixNano() / 1000000
			s.lastTimeStamp = curTimeStamp
			s.sn = 0
		}

		// 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作

		// 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位

		rightBinValue := curTimeStamp & 0x1FFFFFFFFFF

		// 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
		rightBinValue <<= 22

		id := rightBinValue | s.machineID | s.sn
		return id
	}

	if curTimeStamp > s.lastTimeStamp {
		s.sn = 0
		s.lastTimeStamp = curTimeStamp
		s.setLastTimeStamp()

		// 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作

		// 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位

		rightBinValue := curTimeStamp & 0x1FFFFFFFFFF

		// 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
		rightBinValue <<= 22

		id := rightBinValue | s.machineID | s.sn

		return id

	}

	if curTimeStamp < s.lastTimeStamp {
		return 0
	}

	return 0
}

func (s *snowflakeService) GetDisorderID(pid, bid int32) string {
	snowID := s.GetID()
	if snowID == 0 {
		return ""
	}
	timeAt := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s%d%d%d", timeAt, pid, bid, snowID)
}
