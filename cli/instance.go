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

func (c *CLI) instanceCmd() *cobra.Command {
	instanceCmd := cobra.Command{
		Use:   "instance",
		Short: `Manage Dice's service instances`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	return &instanceCmd
}

func (c *CLI) instanceCreateCmd() *cobra.Command {
	var options types.InstanceCreateOptions

	instanceCreateCmd := cobra.Command{
		Use:   "create <SERVICE> <NODE> <URL>",
		Short: `Create a new service instance`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	instanceCreateCmd.Flags().StringVarP(&options.Name, "name", "n", "", `assign a name to the instance`)
	instanceCreateCmd.Flags().StringVarP(&options.Version, "version", "v", "", `specify the deployed service version`)
	instanceCreateCmd.Flags().BoolVarP(&options.Attach, "attach", "a", false, `immediately attach the instance`)

	return &instanceCreateCmd
}

func (c *CLI) instanceAttachCmd() *cobra.Command {
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

func (c *CLI) instanceDetachCmd() *cobra.Command {
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

func (c *CLI) instanceInfoCmd() *cobra.Command {
	var options types.InstanceInfoOptions

	instanceInfoCmd := cobra.Command{
		Use:   "info <ID|NAME|URL>",
		Short: `Print information for a service instance`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}

	instanceInfoCmd.Flags().BoolVarP(&options.Quiet, "quiet", "q", false, `only print the ID`)

	return &instanceInfoCmd
}
