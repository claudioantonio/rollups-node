// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the API for the test DApp.
// It provides functions to encode and decode the test DApp inputs/outputs.
// The inputs/outputs are encoded as binary.
package testdapp

import (
	"bytes"
	"fmt"
)

// Kind of the message.
type MessageKind byte

const (
	// Message that asks the DApp to echo the received payload.
	MessageKindEcho MessageKind = iota

	// Message that asks the DApp to reject the input.
	// The reject message has a payload that should be sent as an report.
	MessageKindReject

	// Number of message kinds.
	numMessages
)

// Get the message kind.
func GetMessageKind(data []byte) (MessageKind, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("empty data")
	}
	kind := MessageKind(data[0])
	if kind < 0 || kind >= numMessages {
		return 0, fmt.Errorf("invalid kind 0x%x", kind)
	}
	return kind, nil
}

// Encode an echo message.
func EncodeEcho(payload []byte) []byte {
	var buf bytes.Buffer
	if err := buf.WriteByte(byte(MessageKindEcho)); err != nil {
		panic(err)
	}
	if _, err := buf.Write(payload); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Decode an echo message.
func DecodeEcho(data []byte) (payload []byte) {
	if MessageKind(data[0]) != MessageKindEcho {
		panic("invalid message kind")
	}
	return data[1:]
}

// Encode a reject message.
func EncodeReject(report []byte) []byte {
	var buf bytes.Buffer
	if err := buf.WriteByte(byte(MessageKindReject)); err != nil {
		panic(err)
	}
	if _, err := buf.Write(report); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// DecodeReject
func DecodeReject(data []byte) (report []byte) {
	if MessageKind(data[0]) != MessageKindReject {
		panic("invalid message kind")
	}
	return data[1:]
}
