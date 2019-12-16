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

// Package client provides the Dice client. While the core package provides
// the daemon, the client is responsible for talking to the daemon's API.
package client

import (
	"github.com/dominikbraun/dice/config"
	"net/http"
)

// setupConfig parses the configuration file and sets all default values
// so that other components can rely on the keys. This step also powers
// Dice's zero-configuration ability.
func (c *Client) setupConfig() error {
	var err error

	if c.config, err = config.NewConfig(configName); err != nil {
		return err
	}

	for key, value := range config.Defaults {
		c.config.SetDefault(key, value)
	}

	return nil
}

// setupInternal sets up the internal HTTP client.
func (c *Client) setupInternal() error {
	c.internal = &http.Client{}
	return nil
}

// setupAPIConnection reads the configured API connection data.
func (c *Client) setupAPIConnection() error {
	c.apiConnection = &APIConnection{
		Protocol: c.config.GetString("api-server-protocol"),
		Host:     c.config.GetString("api-server-host"),
		Port:     c.config.GetString("api-server-port"),
		Root:     c.config.GetString("api-server-root"),
	}

	return nil
}
