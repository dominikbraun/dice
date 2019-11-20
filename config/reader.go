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

// Package config provides utility for parsing configuration values.
package config

import "github.com/spf13/viper"

// Reader prescribes methods for reading configuration values out of files.
// Any type that works such values should exclusively use this interface.
type Reader interface {
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}

// NewReader returns a Reader instance that reads configuration values from
// the specified file. Returns an error of the file cannot be read.
func NewReader(filename string) (Reader, error) {
	reader := viper.New()

	reader.SetConfigName(filename)
	reader.AddConfigPath("../..")

	if err := reader.ReadInConfig(); err != nil {
		return nil, err
	}

	return reader, nil
}
