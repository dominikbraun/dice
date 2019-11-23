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

// newInstanceCommand creates a new command for managing instances.
func newInstanceCommand() *cobra.Command {
	instanceCmd := cobra.Command{
		Use:   "instance",
		Short: `Manage service instances`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	instanceCmd.AddCommand(newInstanceCreateCommand())
	instanceCmd.AddCommand(newInstanceAttachCommand())
	instanceCmd.AddCommand(newInstanceDetachCommand())
	instanceCmd.AddCommand(newInstanceInfoCommand())

	return &instanceCmd
}

// newInstanceCreateCommand creates a new command for creating a new instance.
func newInstanceCreateCommand() *cobra.Command {
	instanceCreateCmd := cobra.Command{
		Use:   "create <SERVICE> <NODE> <URL>",
		Short: `Create a new service instance`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	instanceCreateCmd.Flags().StringP("name", "n", "", `assign a name to the instance`)
	instanceCreateCmd.Flags().StringP("version", "v", "", `specify the deployed service version`)
	instanceCreateCmd.Flags().BoolP("attach", "a", false, `immediately attach the instance`)

	return &instanceCreateCmd
}

// newInstanceAttachCommand creates a new command for attaching an instance.
func newInstanceAttachCommand() *cobra.Command {
	instanceAttachCmd := cobra.Command{
		Use:   "attach <ID|NAME|URL>",
		Short: `Attach an existing service instance`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &instanceAttachCmd
}

// newInstanceDetachCommand creates a new command for detaching an instance.
func newInstanceDetachCommand() *cobra.Command {
	instanceDetachCmd := cobra.Command{
		Use:   "detach <ID|NAME|URL>",
		Short: `Detach an existing service instance`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &instanceDetachCmd
}

// newInstanceInfoCommand creates a new command for printing information.
func newInstanceInfoCommand() *cobra.Command {
	instanceInfoCmd := cobra.Command{
		Use:   "info <ID|NAME|URL>",
		Short: `Print information for a service instance`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	instanceInfoCmd.Flags().BoolP("quiet", "q", false, `only print the ID`)

	return &instanceInfoCmd
}
