package service

import (
    "time"
    "github.com/mfslog/sequenceService/Server/serverPlugin"
    "github.com/mfslog/sequenceService/Server/common"
)


type SnowFlake struct{
    key           string  // etcd key
    machineID     int64 // 机器 id 占10位, 十进制范围是 [ 0, 1023 ]
    sn            int64 // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
    lastTimeStamp int64 // 上次的时间戳(毫秒级), 1秒=1000毫秒, 1毫秒=1000微秒,1微秒=1000纳秒
}

func (sf *SnowFlake)Init() {
    sf.lastTimeStamp = time.Now().UnixNano() / 1000000
    sf.key = common.HostFPath + "/" + "snowflake"
    // 把机器 id 左移 12 位,让出 12 位空间给序列号使用
    sf.machineID = common.MachineID << 12
    sf.getLastTimeStamp()
}

//获得etcd时间
func (sf *SnowFlake)getLastTimeStamp(){
    lastTimeStamp := time.Now().UnixNano() / 1000000
    etcd := serverplugin.GetEtcdConIns();
    sf.lastTimeStamp = etcd.GetInt64Value(sf.key,0)
    if sf.lastTimeStamp == 0{
        sf.lastTimeStamp = lastTimeStamp
        etcd.SetInt64Value(sf.key,sf.lastTimeStamp)
    }
}


//更新etcd 时间记录
func (sf *SnowFlake)setLastTimeStamp(){
    etcd := serverplugin.GetEtcdConIns()
    etcd.SetInt64Value(sf.key, sf.lastTimeStamp)
}


func (sf *SnowFlake)GetSnowflakeId() int64 {
    curTimeStamp := time.Now().UnixNano() / 1000000
    
    // 同一毫秒
    if curTimeStamp == sf.lastTimeStamp {
        sf.sn++
        // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
        if sf.sn > 4095 {
            time.Sleep(time.Millisecond)
            curTimeStamp = time.Now().UnixNano() / 1000000
            sf.lastTimeStamp = curTimeStamp
            sf.sn = 0
        }
        
        // 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作
        
        // 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位
        
        rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
        
        // 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
        rightBinValue <<= 22
        
        id := rightBinValue | sf.machineID | sf.sn
        
        return id
    }
    
    if curTimeStamp > sf.lastTimeStamp {
        sf.sn = 0
        sf.lastTimeStamp = curTimeStamp
        
        // 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作
        
        // 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位
        
        rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
        
        // 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
        rightBinValue <<= 22
        
        id := rightBinValue | sf.machineID | sf.sn
        
        return id
        
    }
    
    if curTimeStamp < sf.lastTimeStamp {
        return 0
    }
    
    return 0
}
