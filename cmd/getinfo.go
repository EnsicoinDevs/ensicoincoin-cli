package cmd

import (
	"context"
	"fmt"
	pb "github.com/EnsicoinDevs/eccctl/rpc"
	"github.com/EnsicoinDevs/eccd/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
	"os"
)

var getinfoCmd = &cobra.Command{
	Use:   "getinfo",
	Short: "Get informations about the node",
	Run: func(cmd *cobra.Command, args []string) {
		client := pb.NewNodeClient(conn)

		reply, err := client.GetInfo(context.Background(), &pb.GetInfoRequest{})
		if err != nil {
			errStatus, _ := status.FromError(err)
			fmt.Println(errStatus.Message())
			os.Exit(1)
		}

		fmt.Printf("\033[1mImplementation:\033[0m %s\n", reply.GetImplementation())
		fmt.Printf("\033[1mProtocol version:\033[0m %d\n", reply.GetProtocolVersion())
		fmt.Printf("\033[1mBest block hash:\033[0m %s\n", utils.NewHash(reply.GetBestBlockHash()).String())
		fmt.Printf("\033[1mGenesis block hash:\033[0m %s\n", utils.NewHash(reply.GetGenesisBlockHash()).String())
	},
}

func init() {
	rootCmd.AddCommand(getinfoCmd)
}
