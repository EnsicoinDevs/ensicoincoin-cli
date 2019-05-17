package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
)

var conn *grpc.ClientConn
var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "eccctl [OPTIONS] <command> <args...>",
	Short: "One cli to rule them all",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rpcserver, err := cmd.Flags().GetString("rpcserver")
		if err != nil {
			fmt.Println("fatal error reading the rpcserver flag", err)
			return
		}

		var opts []grpc.DialOption

		opts = append(opts, grpc.WithInsecure())

		conn, err = grpc.Dial(rpcserver, opts...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("rpcserver", "s", "localhost:4225", "RPC server to connect to")
}
