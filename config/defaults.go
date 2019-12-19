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

// CLIDefaults sets the defaults for CLI-related configuration values.
// They serve as defaults in case the user hasn't specified any other
// values - for the CLI, this can be done with environment variables.
var CLIDefaults = map[string]interface{}{
	"dice-address":     "http://127.0.0.1:9292",
	"dice-api-version": "v1",
}

// DiceDefaults sets the defaults for core-related configuration values.
// They serve as defaults in case the user hasn't specified any other
// values - for the core, this can be done in the Dice config file.
var DiceDefaults = map[string]interface{}{
	"dice-logfile":         "dice.log",
	"api-server-logfile":   "dice.log",
	"proxy-logfile":        "dice.log",
	"kv-store-file":        "dice-store",
	"api-server-protocol":  "http",
	"api-server-host":      "127.0.0.1",
	"api-server-port":      "9292",
	"api-server-root":      "/v1",
	"proxy-port":           "8080",
	"healthcheck-interval": 15000,
	"healthcheck-timeout":  5000,
}
