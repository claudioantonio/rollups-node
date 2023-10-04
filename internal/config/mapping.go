// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package config

import (
	"fmt"
	"os"
)

// This should be removed once we have the whole node rewritten in Go.

func Set() {
	// fills the cache
	getAll()
	GetAuth()

	cache.Lock()
	defer cache.Unlock()

	for external, value := range cache.values {
		internals, ok := mapping[external]
		if !ok {
			panic(fmt.Sprintf("incomplete mapping: (%s, %s)", external, value))
		}

		for _, internal := range internals {
			err := os.Setenv(internal, value)
			if err != nil {
				panic(err)
			}
		}
	}
}

// maps from one external variable to N internal variables
var mapping = map[string][]string{} // TODO
