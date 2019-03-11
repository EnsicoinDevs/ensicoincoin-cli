package cmd

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"os"

	pb "github.com/EnsicoinDevs/ensicoincoin-cli/rpc"
	"github.com/spf13/cobra"
)

var waitforflagCmd = &cobra.Command{
	Use:   "waitforflag",
	Short: "Wait for an accepted tx with a chosen flag",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		flag, _ := cmd.PersistentFlags().GetString("flag")

		var opts []grpc.DialOption

		opts = append(opts, grpc.WithInsecure())

		conn, err := grpc.Dial(os.Getenv("API_NODE"), opts...)
		if err != nil {
			fmt.Println("unable to dial with the grpc server", err)
			return
		}
		defer conn.Close()

		client := pb.NewNodeClient(conn)

		request := &pb.ListenIncomingTxsRequest{}
		stream, err := client.ListenIncomingTxs(context.Background(), request)
		if err != nil {
			fmt.Println("unable to listen for incoming txs", err)
			return
		}

		for {
			txWithBlock, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("unable to receive a tx", err)
				return
			}

			var found bool
			for _, txFlag := range txWithBlock.GetTx().GetFlags() {
				if txFlag == flag {
					found = true
					break
				}
			}

			if found {
				fmt.Println(txWithBlock.GetTx())

				for _, flag := range txWithBlock.GetTx().GetFlags() {
					fmt.Println(flag)
				}

				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(waitforflagCmd)

	waitforflagCmd.PersistentFlags().String("flag", "", "")
}
