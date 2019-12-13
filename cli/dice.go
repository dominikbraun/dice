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

import "github.com/spf13/cobra"

// diceCmd creates and implements the `dice` command, which is also the
// root command. The dice command itself does not have any functionality.
func (c *CLI) diceCmd() *cobra.Command {
	diceCmd := cobra.Command{
		Use:          "dice",
		Short:        `Simple load balancing for non-microservice infrastructures`,
		Long:         `Dice is an ergonomic, flexible, easy to use load balancer designed for non-microservice infrastructures.`,
		Version:      "0.0.0",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &diceCmd
}
