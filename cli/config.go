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

// Package cli provides the Dice CLI commands and their implementation.
package cli

import (
	"errors"
	"github.com/dominikbraun/dice/types"
	"github.com/spf13/cobra"
)

// configCmd creates and implements the `config` command. The config
// command itself does not have any functionality.
func (c *CLI) configCmd() *cobra.Command {
	instanceCmd := cobra.Command{
		Use:   "config",
		Short: `Manage Configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &instanceCmd
}

// configReloadCmd creates and implements the `config reload` command.
func (c *CLI) configReloadCmd() *cobra.Command {
	configReloadCmd := cobra.Command{
		Use:   "reload",
		Short: `Reload configuration`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			route := "/config/reload"

			var response types.Response

			if err := c.client.POST(route, nil, &response); err != nil {
				return err
			}

			if !response.Success {
				return errors.New(response.Message)
			}

			return nil
		},
	}

	return &configReloadCmd
}
