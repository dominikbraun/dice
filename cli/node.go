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
	"github.com/dominikbraun/dice/types"
	"github.com/spf13/cobra"
)

func newNodeCommand() *cobra.Command {
	nodeCmd := cobra.Command{
		Use:   "node",
		Short: `Manage Dice's nodes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	nodeCmd.AddCommand(newNodeCreateCommand())
	nodeCmd.AddCommand(newNodeAttachCommand())
	nodeCmd.AddCommand(newNodeDetachCommand())
	nodeCmd.AddCommand(newNodeInfoCommand())

	return &nodeCmd
}

func newNodeCreateCommand() *cobra.Command {
	var options types.NodeCreateOptions

	nodeCreateCmd := cobra.Command{
		Use:   "create <URL>",
		Short: `Create a new node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	nodeCreateCmd.Flags().StringVarP(&options.Name, "name", "n", "", `assign a name to the node`)
	nodeCreateCmd.Flags().Uint8VarP(&options.Weight, "weight", "w", 1, `specify the node's weight`)
	nodeCreateCmd.Flags().BoolVarP(&options.Attach, "attach", "a", false, `immediately attach the node`)

	return &nodeCreateCmd
}

func newNodeAttachCommand() *cobra.Command {
	nodeAttachCmd := cobra.Command{
		Use:   "attach <ID|NAME|URL>",
		Short: `Attach an existing node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &nodeAttachCmd
}

func newNodeDetachCommand() *cobra.Command {
	nodeDetachCmd := cobra.Command{
		Use:   "detach <ID|NAME|URL>",
		Short: `Detach an existing node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &nodeDetachCmd
}

func newNodeInfoCommand() *cobra.Command {
	var options types.NodeInfoOptions

	nodeInfoCmd := cobra.Command{
		Use:   "info <ID|NAME|URL>",
		Short: `Print information for a node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	nodeInfoCmd.Flags().BoolVarP(&options.Quiet, "quiet", "q", false, `only print the ID`)

	return &nodeInfoCmd
}