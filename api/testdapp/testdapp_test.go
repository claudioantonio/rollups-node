// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package testdapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessageKind(t *testing.T) {
	_, err := GetMessageKind(nil)
	assert.EqualError(t, err, "empty data")

	_, err = GetMessageKind([]byte{0xff})
	assert.EqualError(t, err, "invalid kind 0xff")

	kind, err := GetMessageKind([]byte{0, 1, 2})
	assert.Equal(t, kind, MessageKindEcho)
	assert.Nil(t, err)
}

func TestEncodeDecodeEcho(t *testing.T) {
	data := EncodeEcho([]byte{0xde, 0xad, 0xbe, 0xef})
	assert.Equal(t, []byte{0x00, 0xde, 0xad, 0xbe, 0xef}, data)

	kind, err := GetMessageKind(data)
	assert.Nil(t, err)
	assert.Equal(t, MessageKindEcho, kind)

	payload := DecodeEcho(data)
	assert.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, payload)
}

func TestEncodeDecodeReject(t *testing.T) {
	data := EncodeReject([]byte{0xde, 0xad, 0xbe, 0xef})
	assert.Equal(t, []byte{0x01, 0xde, 0xad, 0xbe, 0xef}, data)

	kind, err := GetMessageKind(data)
	assert.Nil(t, err)
	assert.Equal(t, MessageKindReject, kind)

	report := DecodeReject(data)
	assert.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, report)
}
