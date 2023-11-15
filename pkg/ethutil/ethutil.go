// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains functions to help using the Go-ethereum library.
// It is not the objective of this package to replace or hide Go-ethereum.
package ethutil

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cartesi/rollups-node/pkg/addresses"
	"github.com/cartesi/rollups-node/pkg/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Gas limit when sending transactions.
const GasLimit = 30_000_000

// Dev mnemonic used by Foundry/Anvil.
const FoundryMnemonic = "test test test test test test test test test test test junk"

// Interface that sign blockchain transactions.
type Signer interface {

	// Create the base transaction used in the contract bindings.
	MakeTransactor() (*bind.TransactOpts, error)

	// Get the account address of the signer.
	Account() common.Address
}

// Add input to the input box for the given DApp address.
// This function waits until the transaction is added to a block and return the input index.
func AddInput(
	ctx context.Context,
	client *ethclient.Client,
	addresses *addresses.Book,
	signer Signer,
	input []byte,
) (int, error) {
	inputBox, err := contracts.NewInputBox(addresses.InputBox, client)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to InputBox contract: %v", err)
	}
	receipt, err := sendTransaction(
		ctx, client, signer, big.NewInt(0), GasLimit,
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return inputBox.AddInput(txOpts, addresses.CartesiDApp, input)
		},
	)
	if err != nil {
		return 0, err
	}
	return getInputIndex(ctx, client, addresses, inputBox, receipt)
}

// Get input index in the transaction by looking at the event logs.
func getInputIndex(
	ctx context.Context,
	client *ethclient.Client,
	addresses *addresses.Book,
	inputBox *contracts.InputBox,
	receipt *types.Receipt,
) (int, error) {
	for _, log := range receipt.Logs {
		if log.Address != addresses.InputBox {
			continue
		}
		inputAdded, err := inputBox.ParseInputAdded(*log)
		if err != nil {
			return 0, fmt.Errorf("failed to parse input added event: %v", err)
		}
		// We assume that int will fit all dapp inputs
		inputIndex := int(inputAdded.InputIndex.Int64())
		return inputIndex, nil
	}
	return 0, fmt.Errorf("input index not found")
}

// Get the given input of the given DApp from the input box.
// Return the input sender and the input payload.
func GetInputFromInputBox(
	client *ethclient.Client,
	addresses *addresses.Book,
	inputIndex int,
) (common.Address, []byte, error) {
	inputBox, err := contracts.NewInputBox(addresses.InputBox, client)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to connect to InputBox contract: %v", err)
	}
	it, err := inputBox.FilterInputAdded(
		nil,
		[]common.Address{addresses.CartesiDApp},
		[]*big.Int{big.NewInt(int64(inputIndex))},
	)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to filter input added: %v", err)
	}
	defer it.Close()
	if !it.Next() {
		return common.Address{}, nil, fmt.Errorf("event not found")
	}
	return it.Event.Sender, it.Event.Input, nil
}
