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

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	contentType string = "application/json"
	status404   string = "404 Not Found"
)

var (
	ErrEndpointNotFound = errors.New("the API endpoint could not be found")
)

type Client struct {
	internal    *http.Client
	apiProtocol string
	apiAddress  string
	apiRootPath string
}

func New(apiProtocol, apiAddress, apiRootPath string) *Client {
	c := Client{
		internal:    &http.Client{},
		apiProtocol: apiProtocol,
		apiAddress:  apiAddress,
		apiRootPath: apiRootPath,
	}

	return &c
}

func (c *Client) GET(route string, dest interface{}) error {
	url := c.buildRequestURL(route)

	response, err := c.internal.Get(url)
	if err != nil {
		return err
	}

	if response.Status == status404 {
		return ErrEndpointNotFound
	}

	if err := json.NewDecoder(response.Body).Decode(dest); err != nil && err != io.EOF {
		return err
	}

	_ = response.Body.Close()
	return nil
}

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

	if response.Status == status404 {
		return ErrEndpointNotFound
	}

	if err := json.NewDecoder(response.Body).Decode(dest); err != nil && err != io.EOF {
		return err
	}

	err = response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) buildRequestURL(route string) string {
	rootPath := c.apiRootPath

	if !strings.HasPrefix(rootPath, "/") && len(rootPath) > 0 {
		rootPath = fmt.Sprintf("/%s", rootPath)
	}

	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}

	url := fmt.Sprintf("%s://%s%s%s", c.apiProtocol, c.apiAddress, rootPath, route)
	return url
}
