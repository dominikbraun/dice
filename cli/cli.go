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
	"github.com/dominikbraun/dice/client"
	"github.com/spf13/cobra"
)

// CLI represents the Dice command line interface. It includes all commands
// including arguments and options as well as the HTTP client to call the
// Dice REST API. CLI provides all these commands in its methods.
type CLI struct {
	client  *client.Client
	rootCmd *cobra.Command
}

// New creates a new CLI instance that uses the provided HTTP client.
func New(client *client.Client) *CLI {
	c := CLI{
		client: client,
	}
	c.buildCommands()

	return &c
}

// buildCommands builds all CLI commands by calling the factory functions
// and mounting the child commands to the main commands. All commands have
// to be registered here, otherwise they won't be visible to the user.
func (c *CLI) buildCommands() {
	nodeCmd := c.nodeCmd()

	nodeCmd.AddCommand(c.nodeCreateCmd())
	nodeCmd.AddCommand(c.nodeAttachCmd())
	nodeCmd.AddCommand(c.nodeDetachCmd())
	nodeCmd.AddCommand(c.nodeInfoCmd())
	nodeCmd.AddCommand(c.nodeListCmd())

	serviceCmd := c.serviceCmd()

	serviceCmd.AddCommand(c.serviceCreateCmd())
	serviceCmd.AddCommand(c.serviceEnableCmd())
	serviceCmd.AddCommand(c.serviceDisableCmd())
	serviceCmd.AddCommand(c.serviceInfoCmd())
	serviceCmd.AddCommand(c.serviceListCmd())

	instanceCmd := c.instanceCmd()

	instanceCmd.AddCommand(c.instanceCreateCmd())
	instanceCmd.AddCommand(c.instanceAttachCmd())
	instanceCmd.AddCommand(c.instanceDetachCmd())
	instanceCmd.AddCommand(c.instanceInfoCmd())

	diceCmd := c.diceCmd()

	diceCmd.AddCommand(nodeCmd)
	diceCmd.AddCommand(serviceCmd)
	diceCmd.AddCommand(instanceCmd)

	c.rootCmd = diceCmd
}

// Execute runs the CLI. This means that the command line arguments used
// for running the binary get parsed and processed by cobra.
func (c *CLI) Execute() error {
	return c.rootCmd.Execute()
}
