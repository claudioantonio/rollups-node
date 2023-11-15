// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package send

import (
	"github.com/cartesi/rollups-node/internal/logger"
	"github.com/cartesi/rollups-node/pkg/addresses"
	"github.com/cartesi/rollups-node/pkg/ethutil"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

var (
	ethEndpoint  string
	mnemonic     string
	account      uint32
	hexPayload   string
	envAddresses bool
)

var Cmd = &cobra.Command{
	Use:   "send",
	Short: "Send a rollups input to the Ethereum node",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	logger.Init("info", true)

	payload, err := hexutil.Decode(hexPayload)
	cobra.CheckErr(err)

	ctx := cmd.Context()
	client, err := ethclient.DialContext(ctx, ethEndpoint)
	cobra.CheckErr(err)
	logger.Info.Printf("connected to %v", ethEndpoint)

	signer, err := ethutil.NewMnemonicSigner(ctx, client, mnemonic, account)
	cobra.CheckErr(err)

	var book *addresses.Book
	if envAddresses {
		book = addresses.GetBookFromEnv()
	} else {
		book = addresses.GetTestBook()
	}

	logger.Info.Printf("sending input to %x", book.CartesiDApp)
	inputIndex, err := ethutil.AddInput(ctx, client, book, signer, payload)

	logger.Info.Printf("added input with index %v", inputIndex)
}

func init() {
	Cmd.Flags().StringVar(&ethEndpoint, "eth-endpoint", "http://localhost:8545",
		"ethereum node JSON-RPC endpoint")

	Cmd.Flags().StringVar(&mnemonic, "mnemonic", ethutil.FoundryMnemonic,
		"mnemonic used to sign the transaction")

	Cmd.Flags().Uint32Var(&account, "account", 0,
		"account index used to sign the transaction (default: 0)")

	Cmd.Flags().StringVar(&hexPayload, "payload", "",
		"input payload hex-encoded starting with 0x")
	Cmd.MarkFlagRequired("payload")

	Cmd.Flags().BoolVar(&envAddresses, "env-addresses", false,
		"if set, load contract addresses from env variables; else, use test addresses")
}
