package main

import (
	"fmt"
	"github.com/mfslog/sequenceService/application"
	"github.com/mfslog/sequenceService/domain/snowflake/service"
	"github.com/mfslog/sequenceService/infrastructure/config"
	"github.com/mfslog/sequenceService/infrastructure/driver/etcd"
	"github.com/mfslog/sequenceService/infrastructure/driver/mysql"
	"github.com/mfslog/sequenceService/infrastructure/repository/orderseq_repo"
	"github.com/mfslog/sequenceService/infrastructure/repository/segmentseq_repo"
	"github.com/mfslog/sequenceService/interfaces/rpc"
	sequence "github.com/mfslog/sequenceService/proto"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"net"
	"os"
)

var (
	RootCmd = &cobra.Command{
		Short: "order service",
		Run: func(cmd *cobra.Command, args []string) {
			Run()
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("verison: %d.%d.%d.%d\n", MAJOR, MINOR, PATCH, BUILD)
			os.Exit(0)
		},
	}
)

func buildContainer() *dig.Container {
	container := dig.New()
	container.Provide(config.NewConfig)
	container.Provide(etcd.NewETCDClient)
	container.Provide(mysql.NewMySQLDB)

	container.Provide(orderseq_repo.NewOrderRepo)
	container.Provide(segmentseq_repo.NewSegmentSeqRepo)

	container.Provide(service.NewSnowflakeService)
	container.Provide(application.NewAppService)

	return container
}

func Run() {
	container := buildContainer()

	container.Invoke(func(srv *grpc.Server, app *application.AppService) {
		l , err := net.Listen("tcp",":8080")
		if err != nil{
			fmt.Println(err)
		}
		sequence.RegisterSequenceServer(srv,rpc.Newhandler(app))
		srv.Serve(l)
	})
}

func main() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.Execute()
	return
}
