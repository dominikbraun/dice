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

// Package server provides servers for request proxying and the Dice API.
package server

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

// APIConfig concludes all properties that can be configured by the user.
// Note that the TCP address needs to be secured against remote access.
type APIConfig struct {
	Address string `json:"address"`
	Logfile string `json:"logfile"`
}

// API is an API server which exposes a REST API. The Dice CLI will send all
// requests to this server's endpoints.
type API struct {
	config APIConfig
	router chi.Router
	server *http.Server
}

// NewAPI creates a new API instance and returns a reference to it.
func NewAPI(config APIConfig) *API {
	c := API{
		config: config,
		router: chi.NewRouter(),
	}

	c.server = &http.Server{
		Addr:    c.config.Address,
		Handler: c.router,
	}

	return &c
}

// Run starts the API server. It will listen to the specified port and handle
// incoming requests, sending errors through the returned channel.
//
// When a signal is received through the quit channel, the proxy server attempts
// a graceful shutdown.
func (c *API) Run(quit <-chan os.Signal) chan<- error {
	errorChan := make(chan error)

	go func() {
		err := c.server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			errorChan <- err
		}
	}()

	go func() {
		<-quit

		if err := c.server.Shutdown(context.Background()); err != nil {
			errorChan <- err
		}
	}()

	return errorChan
}
