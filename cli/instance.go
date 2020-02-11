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

// instanceCmd creates and implements the `instance` command. The instance
// command itself does not have any functionality.
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

// instanceCreateCmd creates and implements the `instance create` command.
func (c *CLI) instanceCreateCmd() *cobra.Command {
	var options types.InstanceCreateOptions

	instanceCreateCmd := cobra.Command{
		Use:   "create <SERVICE> <NODE> <URL>",
		Short: `Create a new service instance`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceRef := args[0]
			nodeRef := args[1]
			instanceURL := args[2]
			route := "/instances/create"

			body := types.InstanceCreate{
				ServiceRef:            serviceRef,
				NodeRef:               nodeRef,
				URL:                   instanceURL,
				InstanceCreateOptions: options,
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

	instanceCreateCmd.Flags().StringVarP(&options.Name, "name", "n", "", `assign a name to the instance`)
	instanceCreateCmd.Flags().StringVarP(&options.Version, "version", "v", "", `specify the deployed service version`)
	instanceCreateCmd.Flags().BoolVarP(&options.Attach, "attach", "a", false, `immediately attach the instance`)

	return &instanceCreateCmd
}

// instanceAttachCmd creates and implements the `instance attach` command.
func (c *CLI) instanceAttachCmd() *cobra.Command {
	instanceAttachCmd := cobra.Command{
		Use:   "attach <ID|NAME|URL>",
		Short: `Attach an existing service instance`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			instanceRef := args[0]
			route := "/instances/" + instanceRef + "/attach"

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

	return &instanceAttachCmd
}

// instanceDetachCmd creates and implements the `instance detach` command.
func (c *CLI) instanceDetachCmd() *cobra.Command {
	instanceDetachCmd := cobra.Command{
		Use:   "detach <ID|NAME|URL>",
		Short: `Detach an existing service instance`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			instanceRef := args[0]
			route := "/instances/" + instanceRef + "/detach"

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

	return &instanceDetachCmd
}

// instanceRemoveCmd creates and implemented the `instance remove` command.
func (c *CLI) instanceRemoveCmd() *cobra.Command {
	var options types.InstanceRemoveOptions

	instanceRemoveCmd := cobra.Command{
		Use:     "remove <ID|NAME|URL>",
		Short:   `Remove an instance`,
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			instanceRef := args[0]
			route := "/instances/" + instanceRef + "/remove"

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

	instanceRemoveCmd.Flags().BoolVarP(&options.Force, "force", "f", false, `force the removal`)

	return &instanceRemoveCmd
}

// instanceInfoCmd creates and implements the `instance info` command.
func (c *CLI) instanceInfoCmd() *cobra.Command {
	var options types.InstanceInfoOptions

	instanceInfoCmd := cobra.Command{
		Use:   "info <ID|NAME|URL>",
		Short: `Print information for a service instance`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			instanceRef := args[0]
			route := "/instances/" + instanceRef + "/info"

			var instanceInfoResponse types.InstanceInfoResponse

			if err := c.client.POST(route, nil, &instanceInfoResponse); err != nil {
				return err
			}

			if !instanceInfoResponse.Success {
				return errors.New(instanceInfoResponse.Message)
			}

			fmt.Printf("%v\n", instanceInfoResponse.Data)
			return nil
		},
	}

	instanceInfoCmd.Flags().BoolVarP(&options.Quiet, "quiet", "q", false, `only print the ID`)

	return &instanceInfoCmd
}

// instanceListCmd creates and implements the `instance list` command.
func (c *CLI) instanceListCmd() *cobra.Command {
	var options types.InstanceListOptions

	instanceListCmd := cobra.Command{
		Use:     "list",
		Short:   `List attached instances`,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			route := "/instances/list"
			var instanceListResponse types.InstanceListResponse

			if err := c.client.POST(route, options, &instanceListResponse); err != nil {
				return err
			}

			if !instanceListResponse.Success {
				return errors.New(instanceListResponse.Message)
			}

			for _, n := range instanceListResponse.Data {
				fmt.Printf("%v\n", n)
			}

			return nil
		},
	}

	instanceListCmd.Flags().BoolVarP(&options.All, "all", "a", false, `list all instances`)

	return &instanceListCmd
}
