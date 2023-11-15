// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package manages the contract addresses.
// These addresses usually come from environment variables.
// This package also contain the addresses for the test environment as hard-coded values.
package addresses

import (
	"os"

	"github.com/cartesi/rollups-node/internal/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// List of contract addresses.
type Book struct {
	CartesiDAppFactory  common.Address
	DAppAddressRelay    common.Address
	ERC1155BatchPortal  common.Address
	ERC1155SinglePortal common.Address
	ERC20Portal         common.Address
	ERC721Portal        common.Address
	EtherPortal         common.Address
	InputBox            common.Address
	CartesiDApp         common.Address
}

// Get the addresses for the test environment.
func GetTestBook() *Book {
	return &Book{
		CartesiDAppFactory:  common.HexToAddress("0x7122cd1221C20892234186facfE8615e6743Ab02"),
		DAppAddressRelay:    common.HexToAddress("0xF5DE34d6BbC0446E2a45719E718efEbaaE179daE"),
		ERC1155BatchPortal:  common.HexToAddress("0xedB53860A6B52bbb7561Ad596416ee9965B055Aa"),
		ERC1155SinglePortal: common.HexToAddress("0x7CFB0193Ca87eB6e48056885E026552c3A941FC4"),
		ERC20Portal:         common.HexToAddress("0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB"),
		ERC721Portal:        common.HexToAddress("0x237F8DD094C0e47f4236f12b4Fa01d6Dae89fb87"),
		EtherPortal:         common.HexToAddress("0xFfdbe43d4c855BF7e0f105c400A50857f53AB044"),
		InputBox:            common.HexToAddress("0x59b22D57D4f067708AB0c00552767405926dc768"),
		CartesiDApp:         common.HexToAddress("0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C"),
	}
}

// Get the address book from environment variables.
// Panic if an address is not set.
func GetBookFromEnv() *Book {
	return &Book{
		CartesiDAppFactory:  getAddressFromEnv("CARTESI_DAPP_FACTORY"),
		DAppAddressRelay:    getAddressFromEnv("DAPP_ADDRESS_RELAY"),
		ERC1155BatchPortal:  getAddressFromEnv("ERC1155_BATCH_PORTAL"),
		ERC1155SinglePortal: getAddressFromEnv("ERC1155_SINGLE_PORTAL"),
		ERC20Portal:         getAddressFromEnv("ERC20_PORTAL"),
		ERC721Portal:        getAddressFromEnv("ERC721_PORTAL"),
		EtherPortal:         getAddressFromEnv("ETHER_PORTAL"),
		InputBox:            getAddressFromEnv("INPUT_BOX"),
		CartesiDApp:         getAddressFromEnv("CARTESI_DAPP"),
	}
}

func getAddressFromEnv(varName string) common.Address {
	value := os.Getenv(varName)
	bytes, err := hexutil.Decode(value)
	if err != nil {
		logger.Error.Fatalf("failed to decode address %v", varName)
	}
	if len(bytes) != common.AddressLength {
		logger.Error.Fatalf("address %v with wrong number of bytes", varName)
	}
	return common.Address(bytes)
}
