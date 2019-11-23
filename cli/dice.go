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

// newDiceCommand creates the root command.
func newDiceCommand() *cobra.Command {
	diceCmd := cobra.Command{
		Use:   "dice",
		Short: `Simple load balancing for non-microservice infrastructures`,
		Long:  `Dice is an ergonomic, easy to use load balancer designed for non-microservice infrastructures.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	diceCmd.AddCommand(newNodeCommand())
	diceCmd.AddCommand(newServiceCommand())
	diceCmd.AddCommand(newInstanceCommand())

	return &diceCmd
}

// Build triggers the build of all commands and returns the root command.
func Build() *cobra.Command {
	return newDiceCommand()
}
