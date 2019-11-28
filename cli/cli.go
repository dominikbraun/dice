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

package cli

import (
	"github.com/dominikbraun/dice/client"
	"github.com/spf13/cobra"
)

type CLI struct {
	client  *client.Client
	rootCmd *cobra.Command
}

func New(client *client.Client) *CLI {
	c := CLI{
		client: client,
	}
	c.buildCommands()

	return &c
}

func (c *CLI) buildCommands() {
	nodeCmd := c.nodeCmd()

	nodeCmd.AddCommand(c.nodeCreateCmd())
	nodeCmd.AddCommand(c.nodeAttachCmd())
	nodeCmd.AddCommand(c.nodeDetachCmd())
	nodeCmd.AddCommand(c.nodeInfoCmd())

	serviceCmd := c.serviceCmd()

	serviceCmd.AddCommand(c.serviceCreateCmd())
	serviceCmd.AddCommand(c.serviceEnableCmd())
	serviceCmd.AddCommand(c.serviceDisableCmd())
	serviceCmd.AddCommand(c.serviceInfoCmd())

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

func (c *CLI) Execute() error {
	return c.rootCmd.Execute()
}
