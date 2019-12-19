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

// Package client provides the Dice client. While the core package provides
// the daemon, the client is responsible for talking to the daemon's API.
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dominikbraun/dice/config"
	"io"
	"net/http"
	"strings"
)

const (
	configName  string = "dice"
	contentType string = "application/json"
)

var (
	ErrEndpointNotFound = errors.New("the API endpoint could not be found")
)

// APIConnection stores necessary information for establishing a connection
// to the Dice API server. The values are read from the client's configuration
// reader and can be set via the --address option as well.
type APIConnection struct {
	Address string `json:"address"`
	Version string `json:"root"`
}

// buildURL creates an appropriate URL that can be used to send a request.
func (ac *APIConnection) buildURL() string {
	version := ac.Version

	if !strings.HasPrefix(version, "/") && len(version) > 0 {
		version = fmt.Sprintf("/%s", version)
	}

	url := fmt.Sprintf("%s%s", ac.Address, version)
	return url
}

// Client is the actual Dice client. It is a zero-configuration component
// used by the CLI commands for sending requests and getting responses from
// the API. Configuration values are read every time a command is executed.
type Client struct {
	config        config.Reader
	internal      *http.Client
	apiConnection *APIConnection
}

// New creates a new Client instance and sets up all components.
func New() (*Client, error) {
	var c Client

	if err := c.setup(); err != nil {
		return nil, err
	}

	return &c, nil
}

// setup runs the client setup by invoking all setup* methods.
func (c *Client) setup() error {
	steps := []func() error{
		c.setupConfig,
		c.setupInternal,
		c.setupAPIConnection,
	}

	for _, setup := range steps {
		if err := setup(); err != nil {
			return err
		}
	}

	return nil
}

// OverrideAddress overrides the configured Dice API address permanently.
// This can be useful if a distinct value for the address has been provided,
// for example by using the --address flag of a CLI command.
func (c *Client) OverrideAddress(address string) {
	c.apiConnection.Address = address
}

// GET is the method used by the CLI for sending a GET request to the API.
// If dest is not `nil`, the response body will be decoded into dest.
func (c *Client) GET(route string, dest interface{}) error {
	url := c.buildRequestURL(route)

	response, err := c.internal.Get(url)
	if err != nil {
		return err
	}

	if response.StatusCode == 404 {
		return ErrEndpointNotFound
	}

	if err := json.NewDecoder(response.Body).Decode(dest); err != nil && err != io.EOF {
		return err
	}

	_ = response.Body.Close()
	return nil
}

// POST is the method used by the CLI for sending a POST request to the API.
// If v is not `nil`, it will be encoded into the request body. If dest is
// not `nil`, the response body will be decoded into dest.
func (c *Client) POST(route string, v interface{}, dest interface{}) error {
	url := c.buildRequestURL(route)
	body := bytes.NewBuffer(nil)

	if v != nil {
		if err := json.NewEncoder(body).Encode(v); err != nil {
			return err
		}
	}

	response, err := c.internal.Post(url, contentType, body)
	if err != nil {
		return err
	}

	if response.StatusCode == 404 {
		return ErrEndpointNotFound
	}

	if err := json.NewDecoder(response.Body).Decode(dest); err != nil && err != io.EOF {
		return err
	}

	// Do not handle the error since it is not relevant anymore. The JSON
	// is already decoded and the user is happy.
	_ = response.Body.Close()

	return nil
}

// buildRequestURL creates an entire URL that a request can be sent to. The
// route should be in the form `/my-endpoint`.
func (c *Client) buildRequestURL(route string) string {
	apiURL := c.apiConnection.buildURL()

	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}

	url := fmt.Sprintf("%s%s", apiURL, route)
	return url
}
