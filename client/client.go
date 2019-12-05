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
	"fmt"
	"net/http"
)

const (
	contentType string = "application/json"
)

type Client struct {
	internal   *http.Client
	apiAddress string
}

func New(apiAddress string) *Client {
	c := Client{
		internal:   &http.Client{},
		apiAddress: apiAddress,
	}

	return &c
}

func (c *Client) GET(route string, dest interface{}) error {
	address := fmt.Sprintf("%s%s", c.apiAddress, route)

	response, err := c.internal.Get(address)
	if err != nil {
		return err
	}

	if dest == nil {
		return nil
	}

	return json.NewDecoder(response.Body).Decode(&dest)
}

func (c *Client) POST(route string, v interface{}, dest interface{}) error {
	address := fmt.Sprintf("%s%s", c.apiAddress, route)
	body := new(bytes.Buffer)

	if err := json.NewEncoder(body).Encode(v); err != nil {
		return err
	}

	response, err := c.internal.Post(address, contentType, body)
	if err != nil {
		return err
	}

	if dest == nil {
		return nil
	}

	return json.NewDecoder(response.Body).Decode(&dest)
}
