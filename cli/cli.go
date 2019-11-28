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
	diceCmd *cobra.Command
}

func New(client *client.Client) *CLI {
	c := CLI{
		client: client,
	}
	c.buildCommands()

	return &c
}

func (c *CLI) buildCommands() {
	nodeCmd := c.newNodeCmd()
	nodeCmd.AddCommand(c.newNodeCreateCmd())
	nodeCmd.AddCommand(c.newNodeAttachCmd())
	nodeCmd.AddCommand(c.newNodeDetachCmd())
	nodeCmd.AddCommand(c.newNodeInfoCmd())

	serviceCmd := c.newServiceCmd()
	serviceCmd.AddCommand(c.newServiceCreateCmd())
	serviceCmd.AddCommand(c.newServiceEnableCmd())
	serviceCmd.AddCommand(c.newServiceDisableCmd())
	serviceCmd.AddCommand(c.newServiceInfoCmd())

	instanceCmd := c.newInstanceCmd()
	instanceCmd.AddCommand(c.newInstanceCreateCmd())
	instanceCmd.AddCommand(c.newInstanceAttachCmd())
	instanceCmd.AddCommand(c.newInstanceDetachCmd())
	instanceCmd.AddCommand(c.newInstanceInfoCmd())

	diceCmd := c.newDiceCmd()
	diceCmd.AddCommand(nodeCmd)
	diceCmd.AddCommand(serviceCmd)
	diceCmd.AddCommand(instanceCmd)

	c.diceCmd = diceCmd
}

func (c *CLI) Execute() error {
	return c.diceCmd.Execute()
}
