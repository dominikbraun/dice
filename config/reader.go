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
	"github.com/spf13/viper"
)

// Reader represents a configuration reader. This can be a configuration
// file, system environment variables or other configuration sources.
type Reader interface {
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	SetDefault(key string, value interface{})
}

// NewFile creates a new configuration file reader.
func NewFile(filename string) (Reader, error) {
	r := viper.New()

	r.SetConfigName(filename)
	r.AddConfigPath("/etc/dice/")
	r.AddConfigPath("$HOME/.dice")
	r.AddConfigPath(".")

	if err := r.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return r, nil
		}
		return nil, err
	}

	return r, nil
}

// NewFile creates a new environment variable reader.
func NewEnvironment() (Reader, error) {
	e := Environment{}
	return e, nil
}
