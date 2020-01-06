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
	"fmt"
	"github.com/dominikbraun/dice/types"
	"github.com/spf13/cobra"
)

// nodeCmd creates and implements the `node` command. The node command
// itself does not have any functionality.
func (c *CLI) nodeCmd() *cobra.Command {
	nodeCmd := cobra.Command{
		Use:   "node",
		Short: `Manage Dice's nodes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &nodeCmd
}

// nodeCreateCmd creates and implements the `node create` command.
func (c *CLI) nodeCreateCmd() *cobra.Command {
	var options types.NodeCreateOptions

	nodeCreateCmd := cobra.Command{
		Use:   "create <NAME>",
		Short: `Create a new node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			route := "/nodes/create"

			body := types.NodeCreate{
				Name:              name,
				NodeCreateOptions: options,
			}

			var response types.Response

			if err := c.client.POST(route, body, &response); err != nil {
				return err
			}

			if !response.Success {
				return errors.New(response.Message)
			}

			return nil
		},
	}

	nodeCreateCmd.Flags().Uint8VarP(&options.Weight, "weight", "w", 1, `specify the node's weight`)
	nodeCreateCmd.Flags().BoolVarP(&options.Attach, "attach", "a", false, `immediately attach the node`)

	return &nodeCreateCmd
}

// nodeAttachCmd creates and implements the `node attach` command.
func (c *CLI) nodeAttachCmd() *cobra.Command {
	nodeAttachCmd := cobra.Command{
		Use:   "attach <ID|NAME>",
		Short: `Attach an existing node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeRef := args[0]
			route := "/nodes/" + nodeRef + "/attach"

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

	return &nodeAttachCmd
}

// nodeDetachCmd creates and implements the `node detach` command.
func (c *CLI) nodeDetachCmd() *cobra.Command {
	nodeDetachCmd := cobra.Command{
		Use:   "detach <ID|NAME>",
		Short: `Detach an existing node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeRef := args[0]
			route := "/nodes/" + nodeRef + "/detach"

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

	return &nodeDetachCmd
}

// nodeRemoveCmd creates and implements the `node remove` command.
func (c *CLI) nodeRemoveCmd() *cobra.Command {
	var options types.NodeRemoveOptions

	nodeRemoveCmd := cobra.Command{
		Use:   "remove <ID|NAME>",
		Short: `Remove a node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeRef := args[0]
			route := "/nodes/" + nodeRef + "/remove"

			var response types.Response

			if err := c.client.POST(route, options, &response); err != nil {
				return err
			}

			if !response.Success {
				return errors.New(response.Message)
			}

			return nil
		},
	}

	nodeRemoveCmd.Flags().BoolVarP(&options.Force, "force", "f", false, `force the removal`)

	return &nodeRemoveCmd
}

// nodeInfoCmd creates and implements the `node info` command.
func (c *CLI) nodeInfoCmd() *cobra.Command {
	var options types.NodeInfoOptions

	nodeInfoCmd := cobra.Command{
		Use:   "info <ID|NAME>",
		Short: `Print information for a node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeRef := args[0]
			route := "/nodes/" + nodeRef + "/info"

			var nodeInfoResponse types.NodeInfoResponse

			if err := c.client.POST(route, options, &nodeInfoResponse); err != nil {
				return err
			}

			if !nodeInfoResponse.Success {
				return errors.New(nodeInfoResponse.Message)
			}

			fmt.Printf("%v\n", nodeInfoResponse.Data)
			return nil
		},
	}

	nodeInfoCmd.Flags().BoolVarP(&options.Quiet, "quiet", "q", false, `only print the ID`)

	return &nodeInfoCmd
}

// nodeListCmd creates and implements the `node list` command.
func (c *CLI) nodeListCmd() *cobra.Command {
	var options types.NodeListOptions

	nodeListCmd := cobra.Command{
		Use:   "list",
		Short: `List attached nodes`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			route := "/nodes/list"
			var nodeListResponse types.NodeListResponse

			if err := c.client.POST(route, options, &nodeListResponse); err != nil {
				return err
			}

			if !nodeListResponse.Success {
				return errors.New(nodeListResponse.Message)
			}

			for _, n := range nodeListResponse.Data {
				fmt.Printf("%v\n", n)
			}

			return nil
		},
	}

	nodeListCmd.Flags().BoolVarP(&options.All, "all", "a", false, `list all nodes`)

	return &nodeListCmd
}
