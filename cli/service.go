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
	"github.com/dominikbraun/dice/types"
	"github.com/spf13/cobra"
)

// serviceCmd creates and implements the `service` command. The service
// command itself does not have any functionality.
func (c *CLI) serviceCmd() *cobra.Command {
	serviceCmd := cobra.Command{
		Use:   "service",
		Short: `Manage Dice's services`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &serviceCmd
}

// serviceCreateCmd creates and implements the `service create` command.
func (c *CLI) serviceCreateCmd() *cobra.Command {
	var options types.ServiceCreateOptions

	serviceCreateCmd := cobra.Command{
		Use:   "create <NAME>",
		Short: `Create a new service`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	serviceCreateCmd.Flags().StringVar(&options.Balancing, "balancing", "weighted_round_robin", `specify a balancing method`)
	serviceCreateCmd.Flags().BoolVar(&options.Enable, "enable", false, `immediately enable the service`)

	return &serviceCreateCmd
}

// serviceEnableCmd creates and implements the `service enable` command.
func (c *CLI) serviceEnableCmd() *cobra.Command {
	serviceEnableCmd := cobra.Command{
		Use:   "enable <ID|NAME>",
		Short: `Enable an existing service`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &serviceEnableCmd
}

// serviceDisableCmd creates and implements the `service disable` command.
func (c *CLI) serviceDisableCmd() *cobra.Command {
	serviceDisableCmd := cobra.Command{
		Use:   "disable <ID|NAME>",
		Short: `Disable an existing service`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &serviceDisableCmd
}

// serviceInfoCmd creates and implements the `service info` command.
func (c *CLI) serviceInfoCmd() *cobra.Command {
	var options types.ServiceInfoOptions

	serviceInfoCmd := cobra.Command{
		Use:   "info <ID|NAME>",
		Short: `Print information for a service`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	serviceInfoCmd.Flags().BoolVarP(&options.Quiet, "quiet", "q", false, `only print the ID`)

	return &serviceInfoCmd
}
