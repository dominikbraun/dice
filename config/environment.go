// Copyright 2019 The Dice Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package config provides configuration reader implementations.
package config

import (
	"os"
	"strconv"
)

// Environment represents a set of environment variables. It can be used
// as a config.Reader. Compared to a direct access to environment variables,
// it takes defaults and provides values with different data types.
type Environment map[string]interface{}

// Get implements Reader.Get. Get looks up the environment variable whose
// name is equal to the provided key and returns the value of that variable.
// If it doesn't exist, it searches for a default value. If that fails as
// well, `nil` will be returned.
func (e Environment) Get(key string) interface{} {
	if envVar := os.Getenv(key); envVar != "" {
		return envVar
	}

	if value, ok := e[key]; ok {
		return value
	}

	return nil
}

// GetString implements Reader.GetString. Does the same as Get, but returns
// an empty string if the key cannot be found.
func (e Environment) GetString(key string) string {
	if envVar := os.Getenv(key); envVar != "" {
		return envVar
	}

	if value, ok := e.Get(key).(string); ok {
		return value
	}

	return ""
}

// GetInt implements Reader.GetInt. Does the same as Get, but returns 0
// (zero) if the key cannot be found.
func (e Environment) GetInt(key string) int {
	if envVar := os.Getenv(key); envVar != "" {
		if value, err := strconv.Atoi(envVar); err != nil {
			return value
		}
	}

	if value, ok := e.Get(key).(int); ok {
		return value
	}

	return 0
}

// GetBool implements Reader.GetBool. Does the same as Get, but returns false
// if the key cannot be found.
func (e Environment) GetBool(key string) bool {
	if envVar := os.Getenv(key); envVar != "" {
		if value, err := strconv.ParseBool(envVar); err != nil {
			return value
		}
	}

	if value, ok := e.Get(key).(bool); ok {
		return value
	}

	return false
}

// SetDefault implements Reader.SetDefault. If the default value for the
// specified key already has been set, it will be overridden.
func (e Environment) SetDefault(key string, value interface{}) {
	e[key] = value
}
