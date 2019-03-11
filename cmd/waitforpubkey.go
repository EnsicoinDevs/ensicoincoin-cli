package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"google.golang.org/grpc"
	"io"
	"os"

	pb "github.com/EnsicoinDevs/ensicoincoin-cli/rpc"
	"github.com/spf13/cobra"
)

var waitforpubkeyCmd = &cobra.Command{
	Use:   "waitforpubkey",
	Short: "Wait for an accepted tx with a chosen outpoint",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pubKeyHex, _ := cmd.PersistentFlags().GetString("pubkey")
		pubKeyBytes, _ := hex.DecodeString(pubKeyHex)

		hash := ripemd160.New()
		hash.Write(pubKeyBytes)
		hashHex := hex.EncodeToString(hash.Sum(nil))

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
			for _, output := range txWithBlock.GetTx().GetOutputs() {
				script := output.GetScript()

				if len(script) > 6 {
					script = script[3:]
					script = script[:len(script)-3]

					if hex.EncodeToString(script) == hashHex {
						found = true
						break
					}
				}
			}

			if found {
				fmt.Println(txWithBlock.GetTx().Hash)

				for _, flag := range txWithBlock.GetTx().GetFlags() {
					fmt.Println(flag)
				}

				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(waitforpubkeyCmd)

	waitforpubkeyCmd.PersistentFlags().String("pubkey", "", "")
}
