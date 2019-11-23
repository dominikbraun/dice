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

// newServiceCommand creates a new command for managing services.
func newServiceCommand() *cobra.Command {
	serviceCmd := cobra.Command{
		Use:   "service",
		Short: `Manage Dice's services`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	serviceCmd.AddCommand(newServiceCreateCommand())
	serviceCmd.AddCommand(newServiceActivateCommand())
	serviceCmd.AddCommand(newServiceDeactivateCommand())
	serviceCmd.AddCommand(newServiceInfoCommand())

	return &serviceCmd
}

// newServiceCreateCommand creates a new command for creating a service.
func newServiceCreateCommand() *cobra.Command {
	serviceCreateCmd := cobra.Command{
		Use:   "create <NAME>",
		Short: `Create a new service`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	serviceCreateCmd.Flags().String("balancing", "weighted_round_robin", `specify a balancing method`)
	serviceCreateCmd.Flags().Bool("activate", false, `immediately activate the service`)

	return &serviceCreateCmd
}

// newServiceActivateCommand creates a new command for activating a service.
func newServiceActivateCommand() *cobra.Command {
	serviceActivateCmd := cobra.Command{
		Use:   "activate <ID|NAME>",
		Short: `Activate an existing service`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &serviceActivateCmd
}

// newServiceDeactivateCommand creates a new command for activating a service.
func newServiceDeactivateCommand() *cobra.Command {
	serviceDeactivateCmd := cobra.Command{
		Use:   "deactivate <ID|NAME>",
		Short: `Deactivate an existing service`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &serviceDeactivateCmd
}

// newServiceInfoCommand creates a new command for printing information.
func newServiceInfoCommand() *cobra.Command {
	serviceInfoCmd := cobra.Command{
		Use:   "info <ID|NAME>",
		Short: `Print information for a service`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	serviceInfoCmd.Flags().BoolP("quiet", "q", false, `only print the ID`)

	return &serviceInfoCmd
}
