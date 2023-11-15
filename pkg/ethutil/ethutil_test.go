// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package ethutil

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/cartesi/rollups-node/pkg/addresses"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		err := anvilContainer.Terminate(ctx)
		assert.Nil(t, err)
	}()

	endpoint, err := anvilContainer.Endpoint(ctx, "ws")
	require.Nil(t, err)

	client, err := ethclient.DialContext(ctx, endpoint)
	require.Nil(t, err)

	signer, err := NewMnemonicSigner(ctx, client, FoundryMnemonic, 0)
	require.Nil(t, err)

	book := addresses.GetTestBook()
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
				return AddInput(ctx, client, book, signer, payload)
			},
			sender: common.HexToAddress("f39fd6e51aad88f6f4ce6ab8827279cfffb92266"),
			input:  payload,
		},
	}
	for i, testCase := range testCases {
		t.Logf("testing client.%v", testCase.name)

		inputIndex, err := testCase.do()
		if !assert.Nil(t, err) {
			logContainerOutput(t, ctx, anvilContainer)
			t.FailNow()
		}

		require.Equal(t, i, inputIndex)

		event, err := GetInputFromInputBox(client, book, inputIndex)
		require.Nil(t, err)
		require.Equal(t, testCase.sender, event.Sender)
		require.Equal(t, testCase.input, event.Input)
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
	require.Nil(t, err)
	return anvilContainer
}

// Log the output of the given container
func logContainerOutput(t *testing.T, ctx context.Context, container testcontainers.Container) {
	reader, err := container.Logs(ctx)
	require.Nil(t, err)
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	require.Nil(t, err)
	t.Log(string(bytes))
}
