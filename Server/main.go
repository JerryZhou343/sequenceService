package main

import (
    "path/filepath"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "fmt"
    "net"
    configer2 "github.com/mfslog/sequenceService/Server/configer"
    "github.com/mfslog/sequenceService/Server/common"
    "github.com/mfslog/sequenceService/Server/log"
    "github.com/mfslog/sequenceService/Server/DBSession"
    "github.com/mfslog/sequenceService/Server/serverPlugin"
    "github.com/mfslog/sequenceService/Server/service"
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
            
            //2.配置日志信息
            log.SetupLogger();
            
            //3.加载etcd配置内容
            configer := configer2.GetInstance()
            configer.LocadConfig()
            
            //4.连接mysql数据库
            dbInstance := DBSession.GetInstance()
            dbInstance.InitDBCon()
            
            
            //5.连接redis 数据库
            //cacheInstance := cacheSession.GetInstance()
            //cacheInstance.InitRedisCon()
            
            //6.连接kafka
            
            //kafka := kafka2.GetProducerInstance()
            //kafka.Init()
            
            //6.向etcd 注册服务
            // 获取IP 和  port
            ifaces, _ := net.Interfaces()
            for _, i := range  ifaces{
                addrs, _ := i.Addrs()
                for _, addr := range addrs{
                    switch v := addr.(type) {
                    case *net.IPAddr:
                        common.IP = v.IP.String()
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
            
            //7.监听端口,对外服务
           service.NewServer(common.Port)
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
    //断开kafka 连接
}