package main

import (
    "path/filepath"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "fmt"
    configer2 "github.com/mfslog/sequenceService/Server/configer"
    "github.com/mfslog/sequenceService/Server/common"
    "github.com/mfslog/sequenceService/Server/log"
    "github.com/mfslog/sequenceService/Server/DBSession"
    "github.com/mfslog/sequenceService/Server/serverPlugin"
    "github.com/mfslog/sequenceService/Server/cacheSession"
    "net"
    "strings"
    "github.com/mfslog/sequenceService/Server/transport"
    "github.com/oklog/run"
    "github.com/go-kit/kit/metrics/prometheus"
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/go-kit/kit/metrics"
    stdprometheus "github.com/prometheus/client_golang/prometheus"
    "google.golang.org/grpc"
    pb "github.com/mfslog/sequenceService/proto"

    stdzipkin "github.com/openzipkin/zipkin-go"
    zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func main(){
    
    var configFile string
    common.ApplicationName = filepath.Base(os.Args[0])
    rootCmd := &cobra.Command{
        Use:"",
        Short:common.ApplicationName +" [flag]",
        Long:common.ApplicationName + " [flag",
        Run: func(cmd *cobra.Command, args []string) {
            
            //1.加载本地配置文件
            viper.SetConfigFile(configFile)
            viper.SetConfigType("json")
            err := viper.ReadInConfig()
            configFound := true
            if err != nil {
                switch err.(type) {
                case viper.ConfigParseError:
                    fmt.Println("Error parsing configuration: %s\n", err)
                    os.Exit(1)
                default:
                    configFound = false
                }
            }
            if !configFound {
                fmt.Println("can't find configer file")
                os.Exit(1)
            }

            common.MachineID = viper.GetInt64("common.machine_id")
            basePath := viper.GetString("service_register.base_path")
            name := viper.GetString("service_register.service_name")
            common.LockPath = strings.TrimRight(basePath,"/") + "/" + name + "/dblock"

            //2.配置日志信息
            log.SetupLogger();
            
            //3.加载etcd配置内容
            etcd := serverplugin.GetEtcdConIns()
            etcd.Init()
            configer := configer2.GetInstance()
            configer.LocadConfig()
            
            //4.连接mysql数据库
            dbInstance := DBSession.GetInstance()
            dbInstance.InitDBCon()
            
            
            //5.连接redis 数据库
            cacheInstance := cacheSession.GetInstance()
            cacheInstance.InitCon()
            

            
            //6.向etcd 注册服务
            // 获取IP 和  port
           addrs, err := net.InterfaceAddrs()
           if err != nil{
               log.Error("can't get host ip")
               os.Exit(1)
           }

           for _,address := range addrs{
               if ipnet, ok := address.(*net.IPNet); ok &&
                   !ipnet.IP.IsLoopback() &&
                   !ipnet.IP.IsMulticast() {
                   if ipnet.IP.To4() != nil{
                       common.IP = ipnet.IP.String()
                   }
               }
           }

            common.Port = viper.GetInt("service_register.service_port")
            //etcd service
            serviceInfo := serverplugin.ServiceInfo{
                Name: viper.GetString("service_register.service_name"),
                IP: common.IP,
                Port: viper.GetInt("service_register.service_port"),
                BasePath:viper.GetString("service_register.base_path"),
            }
            
            //load configer from etcd
            etcdService , _ := serverplugin.NewEtcdService(serviceInfo,viper.GetStringSlice("service_register.etcd_address"))
            etcdService.RegisterService()
            var CallSeqCounter metrics.Counter
            var zipkinTracer *stdzipkin.Tracer
            {
                var (
                    err error
                    hostPort = "0.0.0.0:8486"
                    serviceName = "sequence"
                    useNoopTracer = false
                    reporter = zipkinhttp.NewReporter(viper.GetString("common.zipkinUrl"))
                )
                defer reporter.Close()
                zEP , _ := stdzipkin.NewEndpoint(serviceName,hostPort)
                zipkinTracer, err = stdzipkin.NewTracer(
                    reporter,stdzipkin.WithLocalEndpoint(zEP),stdzipkin.WithNoopTracer(useNoopTracer),
                )
                if err != nil{
                    os.Exit(1)
                }
            }
            var g run.Group
            {
                CallSeqCounter  = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
                    Namespace: "private",
                    Subsystem: "sequence",
                    Name:      "call_summed",
                    Help:      "Total count of get sequence via the getSequence method.",
                }, []string{})

                http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
                addr := fmt.Sprint("0.0.0.0:",viper.GetInt("service_register.metric_port"))
                fmt.Println("test metric",addr)
                metricListener,err := net.Listen("tcp",addr)
                if err != nil{
                    os.Exit(1)
                }
                g.Add(func()error{

                    return http.Serve(metricListener,http.DefaultServeMux)
                }, func(error){
                    metricListener.Close()
                })
            }
            //7.监听端口,对外服务
            {

                serviceAddress := fmt.Sprintf("0.0.0.0:%d",common.Port)
                log.Info("listen:" + serviceAddress)

                ls, _ := net.Listen("tcp", serviceAddress)
                g.Add(func() error {
                    baseServer := grpc.NewServer()
                    pb.RegisterSequenceServer(baseServer, transport.NewGrpcServer(CallSeqCounter,zipkinTracer))
                    return baseServer.Serve(ls)
                }, func(e error) {
                   ls.Close()
                })
            }
            g.Run()
            uninitialize()
        },
    }
    
    
    rootCmd.Flags().StringVarP(&configFile,"config","c","config.json","config file")
    
    versionCmd := &cobra.Command{
        Use:"version",
        Short:"",
        Long:"",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println(common.GetVersionInfo())
        },
    }
    
    rootCmd.AddCommand(versionCmd)
    rootCmd.Execute()
}


func uninitialize(){
    //断开mysql连接
   dbIns := DBSession.GetInstance()
   dbIns.UninitDBCon()

   redisIns := cacheSession.GetInstance()
   redisIns.Close()


   etcdIns := serverplugin.GetEtcdConIns()
   etcdIns.Close()
}