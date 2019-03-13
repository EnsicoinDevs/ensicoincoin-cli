package cmd

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"google.golang.org/grpc"
	"os"
	"time"

	pb "github.com/EnsicoinDevs/ensicoincoin-cli/rpc"
	"github.com/EnsicoinDevs/ensicoincoin/network"
	"github.com/EnsicoinDevs/ensicoincoin/scripts"
	"github.com/EnsicoinDevs/ensicoincoin/utils"
	"github.com/btcsuite/btcd/btcec"
	"github.com/spf13/cobra"
)

// sendtoCmd represents the sendto command
var sendtoCmd = &cobra.Command{
	Use:   "sendto",
	Short: "Create and publish a transaction",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		tx := network.NewTxMessage()

		hash, _ := cmd.PersistentFlags().GetString("outpointhash")
		index, _ := cmd.PersistentFlags().GetUint32("outpointindex")

		privKeyHex, _ := cmd.PersistentFlags().GetString("privkey")
		pubKeyHex, _ := cmd.PersistentFlags().GetString("pubkey")
		spentOutputValue, _ := cmd.PersistentFlags().GetUint64("spentoutputvalue")

		hashBytes, _ := hex.DecodeString(hash)

		tx.Flags, _ = cmd.PersistentFlags().GetStringArray("flag")

		tx.Inputs = []*network.TxIn{
			&network.TxIn{
				PreviousOutput: &network.Outpoint{*utils.NewHash(hashBytes), index},
			},
		}

		value, _ := cmd.PersistentFlags().GetUint64("value")

		tx.Outputs = []*network.TxOut{
			&network.TxOut{
				Value:  value,
				Script: []byte{byte(scripts.OP_TRUE)},
			},
		}

		if pubKeyHex != "" {
			// OP_DUP OP_HASH160

			tx.Outputs[0].Script = []byte{byte(scripts.OP_DUP), byte(scripts.OP_HASH160)}

			pubKeyBytes, _ := hex.DecodeString(pubKeyHex)

			hash := ripemd160.New()
			hash.Write(pubKeyBytes)
			pubKeyHashBytes := hash.Sum(nil)
			pubKeyHashBytesSize := len(pubKeyHashBytes)

			// <hash160(pubKey)>

			tx.Outputs[0].Script = append(tx.Outputs[0].Script, byte(pubKeyHashBytesSize))
			tx.Outputs[0].Script = append(tx.Outputs[0].Script, pubKeyHashBytes...)

			// OP_EQUAL OP_VERIFY OP_CHECKSIG

			tx.Outputs[0].Script = append(tx.Outputs[0].Script, []byte{byte(scripts.OP_EQUAL), byte(scripts.OP_VERIFY), byte(scripts.OP_CHECKSIG)}...)
		}

		if privKeyHex != "" {
			shash := tx.SHash(tx.Inputs[0], spentOutputValue)

			privKeyBytes, _ := hex.DecodeString(privKeyHex)

			privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)

			sig, err := privKey.Sign(shash[:])
			if err != nil {
				fmt.Println("error signing the tx", err)
				return
			}

			sigBytes := sig.Serialize()
			sigBytesSize := len(sigBytes)

			// <sig>

			tx.Inputs[0].Script = append([]byte{byte(sigBytesSize)}, sigBytes...)

			pubKeyBytes := pubKey.SerializeCompressed()
			pubKeyBytesSize := len(pubKeyBytes)

			// <pubKey>

			tx.Inputs[0].Script = append(tx.Inputs[0].Script, append([]byte{byte(pubKeyBytesSize)}, pubKeyBytes...)...)
		}

		buf := bytes.NewBuffer(nil)
		tx.Encode(buf)

		var opts []grpc.DialOption

		opts = append(opts, grpc.WithInsecure())

		conn, err := grpc.Dial(os.Getenv("API_NODE"), opts...)
		if err != nil {
			fmt.Println("unable to dial with the grpc server")
			return
		}
		defer conn.Close()

		client := pb.NewNodeClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = client.PublishTx(ctx, &pb.PublishTxRequest{
			Tx: buf.Bytes(),
		})
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(hex.EncodeToString(tx.Hash()[:]))
	},
}

func init() {
	rootCmd.AddCommand(sendtoCmd)

	sendtoCmd.PersistentFlags().String("outpointhash", "", "")
	sendtoCmd.PersistentFlags().Uint32("outpointindex", 0, "")
	sendtoCmd.PersistentFlags().Uint64("value", 0, "")
	sendtoCmd.PersistentFlags().String("pubkey", "", "")
	sendtoCmd.PersistentFlags().String("privkey", "", "")
	sendtoCmd.PersistentFlags().Uint64("spentoutputvalue", 0, "")
	sendtoCmd.PersistentFlags().StringArray("flag", nil, "")
}
