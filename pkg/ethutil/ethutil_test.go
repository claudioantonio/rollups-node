// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package ethutil

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const testTimeout = 300 * time.Second

// Test sending inputs to the input box.
// This function test the GetInputFromInputBox function as well.
func TestSendingInputs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	anvilContainer := setupEthContainer(t, ctx)
	defer func() {
		if err := anvilContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate anvil container: %v", err)
		}
	}()

	endpoint, err := anvilContainer.Endpoint(ctx, "ws")
	if err != nil {
		t.Fatalf("failed to get anvil endpoint: %v", err)
	}

	client, err := ethclient.DialContext(ctx, endpoint)
	if err != nil {
		t.Fatalf("failed to connect to eth node: %v", err)
	}

	signer, err := NewMnemonicSigner(ctx, client, FoundryMnemonic, 0)
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}

	dappAddress := common.HexToAddress("fafafafafafafafafafafafafafafafafafafafa")
	payload := common.Hex2Bytes("deadbeef")

	testCases := []struct {
		name   string
		do     func() (int, error)
		sender common.Address
		input  []byte
	}{
		{
			name: "AddInput",
			do: func() (int, error) {
				return AddInput(ctx, client, signer, dappAddress, payload)
			},
			sender: common.HexToAddress("f39fd6e51aad88f6f4ce6ab8827279cfffb92266"),
			input:  payload,
		},
	}
	for i, testCase := range testCases {
		t.Logf("testing client.%v", testCase.name)
		inputIndex, err := testCase.do()
		if err != nil {
			logContainerOutput(t, ctx, anvilContainer)
			t.Fatalf("failed to send: %v", err)
		}
		if inputIndex != i {
			t.Fatalf("wrong input index: %v; expected: %v", inputIndex, i)
		}
		readSender, readInput, err := GetInputFromInputBox(client, dappAddress, inputIndex)
		if err != nil {
			t.Fatal(err)
		}
		if readSender != testCase.sender {
			t.Fatalf("wrong sender: %x", readSender)
		}
		if !bytes.Equal(readInput, testCase.input) {
			t.Fatalf("wrong input: %x", readInput)
		}
	}
}

// We use the sunodo devnet docker image to test the client.
// This image starts an anvil node with the Rollups contracts already deployed.
func setupEthContainer(t *testing.T, ctx context.Context) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image: "sunodo/devnet:1.1.1",
		Cmd: []string{
			"anvil",
			"--block-time",
			"1",
			"--load-state",
			"/usr/share/sunodo/anvil_state.json",
		},
		Env: map[string]string{
			"ANVIL_IP_ADDR": "0.0.0.0",
		},
		ExposedPorts: []string{"8545/tcp"},
		WaitingFor:   wait.ForLog("Listening on 0.0.0.0:8545"),
	}
	anvilContainer, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		t.Fatalf("failed to start anvil container: %v", err)
	}
	return anvilContainer
}

func logContainerOutput(t *testing.T, ctx context.Context, container testcontainers.Container) {
	reader, err := container.Logs(ctx)
	if err != nil {
		t.Fatalf("failed to get reader: %v", err)
	}
	bytes, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read logs: %v", err)
	}
	t.Log(string(bytes))
}
