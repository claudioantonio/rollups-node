// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

//go:generate go run generate/main.go

// Parsing functions.
var (
	toBool     = strconv.ParseBool
	toInt      = strconv.Atoi
	toInt64    = func(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }
	toString   = func(s string) (string, error) { return s, nil }
	toDuration = func(s string) (time.Duration, error) { return time.ParseDuration(s + "s") }
	toLogLevel = func(s string) (LogLevel, error) {
		var logLevel = map[string]LogLevel{
			"debug":   LogLevelDebug,
			"info":    LogLevelInfo,
			"warning": LogLevelWarning,
			"error":   LogLevelError,
		}
		if ll, ok := logLevel[s]; ok {
			return ll, nil
		} else {
			return LogLevelDebug, fmt.Errorf(`invalid log level "%s"`, s)
		}
	}
)

func GetServerManagerSessionId() string {
	return "default_session_id"
}

func GetInspectServerEndpoint() string {
	port := getInspectServerPort()
	if port == nil {
		panic("")
	}
	return fmt.Sprintf("0.0.0.0:%d", *port)
}

func GetAuth() Auth {
	// getting the (optional) account index
	index := getTxSigningMnemonicAccountIndex()

	// if the mnemonic is coming from an environment variable
	if mnemonic := getTxSigningMnemonic(); mnemonic != nil {
		return AuthMnemonic{*mnemonic, index}
	}

	// if the mnemonic is coming from a file
	if file := getTxSigningMnemonicFile(); file != nil {
		mnemonic, err := os.ReadFile(*file)
		if err != nil {
			panic(err)
		}
		return AuthMnemonic{string(mnemonic), index}
	}

	// if we are not using mnemonics, but AWS authentication
	keyID := getTxSigningAwsKmsKeyId()
	region := getTxSigningAwsKmsRegion()
	if keyID == nil || region == nil {
		panic("missing tx-signing environment variables")
	}
	return AuthAWS{*keyID, *region}
}

// ------------------------------------------------------------------------------------------------
// Custom Types
// ------------------------------------------------------------------------------------------------

type LogLevel uint8

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

// ------------------------------------------------------------------------------------------------

type Auth interface{}

type AuthMnemonic struct {
	Mnemonic     string
	AccountIndex *int
}

type AuthAWS struct {
	KeyID  string
	Region string
}

// ------------------------------------------------------------------------------------------------
// Get Helpers
// ------------------------------------------------------------------------------------------------

const prefix = "CARTESI_"

// Cache of environment variable values
var cache struct {
	sync.Mutex
	values map[string]string
}

// Reads the value of an environment variable (loads from a cached value when possible).
// It returns the value read and true if the variable was set,
// otherwise it returns the empty string and false.
func read(name string) (string, bool) {
	cache.Lock()
	defer cache.Unlock()
	if s, ok := cache.values[name]; ok {
		return s, true
	} else if s, ok := os.LookupEnv(prefix + name); ok {
		cache.values[name] = s
		return s, true
	} else {
		return "", false
	}
}

// Parses a string using a given function.
// It panics on a parsing error, otherwise returns the parsed value.
func parse[T any](s string, f func(string) (T, error)) T {
	v, err := f(s)
	if err != nil {
		panic(err)
	}
	return v
}

// Returns nil or the value of an environment variable.
//
// If the variable could not be read from the environment and it has a default value,
// then this function will set the cache with the default value and return it parsed.
func getOptional[T any](name, default_ string, hasDefault bool, parser func(string) (T, error)) *T {
	if s, ok := read(name); ok {
		v := parse(s, parser)
		return &v
	}

	if hasDefault {
		cache.Lock()
		defer cache.Unlock()
		cache.values[name] = default_
		v := parse(default_, parser)
		return &v
	}

	return nil
}

// Same as getOptional, but panics instead of returning nil.
func get[T any](name, defaultValue string, hasDefault bool, parser func(string) (T, error)) T {
	v := getOptional(name, defaultValue, hasDefault, parser)
	if v == nil {
		panic(fmt.Errorf(`missing required %s env var`, prefix+name))
	}

	return *v
}
