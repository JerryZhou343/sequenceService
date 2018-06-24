package log

import (
    "github.com/go-kit/kit/log"
    "fmt"
    "github.com/spf13/viper"
    "github.com/go-kit/kit/log/level"
    //"bytes"
    "os"
    "github.com/mfslog/sequenceService/Server/common"
    "strings"
    "strconv"
)

var(
    logger log.Logger
)



//type log2kafka struct{
//    client *kafka.KafkaProducerClient
//}
//
//func (k *log2kafka)init(){
//    k.client = kafka.GetProducerInstance()
//}
//
//func (k *log2kafka)Write(p []byte) (n int, err error){
//    msg := &sarama.ProducerMessage{
//        Topic: "log",
//        Value: sarama.ByteEncoder(p),
//    }
//    k.client.Producer.Input() <- msg
//    return len(p),nil
//}
//


/// Caller returns a Valuer that returns a file and line from a specified depth
// in the callstack. Users will probably want to use DefaultCaller.
var LogCaller = log.Caller(6)

func SetupLogger(){
    logPath := viper.GetString("common.log_path")
    _,err := os.Open(logPath)
    if err != nil{
        err := os.MkdirAll(logPath,os.ModePerm)
        if err != nil{
            fmt.Println("can't mkdir" + logPath)
        }
    }
    
    absFile := strings.TrimRight(logPath, "/") + "/" + common.ApplicationName + "_m" + strconv.FormatInt(common.MachineID,10)  + ".log"
    //var fd *log2kafka = &log2kafka{};
    //fd.init();
    var fd *os.File
    if common.CheckFileIsExist(absFile) { //如果文件存在
        fd, _ = os.OpenFile(absFile, os.O_APPEND, 0666) //打开文件
        //fmt.Println("文件存在")
    } else {
        fd, _ = os.Create(absFile) //创建文件
        //fmt.Println("文件不存在")
    }
    
    //logger.Log(common.GetVersionInfo())
    
    logger = log.NewJSONLogger(log.NewSyncWriter(fd))
    logger = log.With(logger,
        //"module", MineDodule,
        "ts", log.DefaultTimestamp,
        "caller", LogCaller,
    )
    logger = level.NewFilter(logger, level.AllowAll())
}


func Info(msg string){
    level.Info(logger).Log("module", viper.GetString("service_register.service_name"),"msg",msg)
}

func Error(msg string){
    level.Error(logger).Log("module",viper.GetString("service_register.service_name"), "msg",msg)
}

func Debug(msg string){
    level.Debug(logger).Log("module",viper.GetString("service_register.service_name"),"msg",msg)
}