package main

import "github.com/spf13/cobra"

var (
	RootCmd = &cobra.Command{
		Use:"seq",
		Run: func(cmd *cobra.Command, args []string) {
			Run()
		},
	}


)

func Run(){


}


func main(){
	RootCmd.Execute()
}
