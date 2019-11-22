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

// Package cli provides all CLI command and the Dice API client.
package cli

import "github.com/spf13/cobra"

// newNodeCommand creates a new command for managing nodes.
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
	nodeCmd.AddCommand(newNodeInfoCommand())

	return &nodeCmd
}

// newNodeCreateCommand creates a command for creating a new node.
func newNodeCreateCommand() *cobra.Command {
	nodeCreateCmd := cobra.Command{
		Use:   "create <URL>",
		Short: `Create a new node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	nodeCreateCmd.Flags().StringP("name", "n", "", `assign a name to the node`)
	nodeCreateCmd.Flags().Uint8P("weight", "w", 1, `specify the node's weight`)
	nodeCreateCmd.Flags().BoolP("attach", "a", false, `immediately attach the node`)

	return &nodeCreateCmd
}

// newNodeAttachCommand creates a new command for attaching a node.
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

// newNodeInfoCommand creates a new command for printing information.
func newNodeInfoCommand() *cobra.Command {
	nodeInfoCmd := cobra.Command{
		Use:   "info <ID|NAME|URL>",
		Short: `Print information for a node`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	nodeInfoCmd.Flags().BoolP("quiet", "q", false, `only print the ID`)

	return &nodeInfoCmd
}
