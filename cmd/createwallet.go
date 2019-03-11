package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/spf13/cobra"
)

var createwalletCmd = &cobra.Command{
	Use:   "createwallet",
	Short: "Create a new wallet",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		privateKey, err := btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			fmt.Println("error generating a new private key", err)
			return
		}

		publicKey := privateKey.PubKey()

		fmt.Println("privateKey:", hex.EncodeToString(privateKey.Serialize()))
		fmt.Println("publicKey:", hex.EncodeToString(publicKey.SerializeCompressed()))
	},
}

func init() {
	rootCmd.AddCommand(createwalletCmd)
}
