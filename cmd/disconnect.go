package cmd

import (
	"context"
	"fmt"
	pb "github.com/EnsicoinDevs/eccctl/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"strconv"
)

var disconnectCmd = &cobra.Command{
	Use:   "disconnect [address]",
	Short: "Disconnect the node from a peer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := pb.NewNodeClient(conn)

		host, port, err := net.SplitHostPort(args[0])
		if err != nil {
			fmt.Printf("invalid address: %s\n", err)
			os.Exit(1)
		}

		ip, err := net.LookupIP(host)
		if err != nil {
			fmt.Printf("invalid address: %s\n", err)
			os.Exit(1)
		}

		parsedPort, err := strconv.Atoi(port)
		if err != nil {
			fmt.Printf("invalid address: %s\n", err)
			os.Exit(1)
		}

		request := &pb.DisconnectPeerRequest{
			Peer: &pb.Peer{
				Address: &pb.Address{
					Ip:   ip[0].String(),
					Port: uint32(parsedPort),
				},
			},
		}
		reply, err := client.DisconnectPeer(context.Background(), request)
		if err != nil {
			errStatus, _ := status.FromError(err)
			fmt.Println(errStatus.Message())
			os.Exit(1)
		}

		_ = reply
	},
}

func init() {
	rootCmd.AddCommand(disconnectCmd)
}
