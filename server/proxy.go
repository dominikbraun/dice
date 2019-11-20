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
	"github.com/dominikbraun/dice/scheduler"
	"net/http"
	"os"
)

// ProxyConfig concludes all properties that can be configured by the user.
type ProxyConfig struct {
	Address string `json:"address"`
	Logfile string `json:"logfile"`
}

// Proxy is a proxy server that will handle incoming requests, establish a
// new connection to an service instance and return the service's respond.
type Proxy struct {
	config    ProxyConfig
	server    *http.Server
	scheduler scheduler.Scheduler
}

// NewProxy creates a new Proxy instance and returns a reference to it.
func NewProxy(config ProxyConfig, scheduler scheduler.Scheduler) *Proxy {
	p := Proxy{
		config:    config,
		scheduler: scheduler,
	}

	p.server = &http.Server{
		Addr:    p.config.Address,
		Handler: p.handleRequest(),
	}

	return &p
}

// Run starts the proxy server. It will listen to the specified port and
// handle incoming requests, sending errors through the returned channel.
//
// When a signal is received through the quit channel, the proxy server
// attempts a graceful shutdown.
func (p *Proxy) Run(quit <-chan os.Signal) chan<- error {
	errorChan := make(chan error)

	go func() {
		err := p.server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			errorChan <- err
		}
	}()

	go func() {
		<-quit

		if err := p.server.Shutdown(context.Background()); err != nil {
			errorChan <- err
		}
	}()

	return errorChan
}

// handleRequest implements the core request handling logic which includes
// the transfer/copy of the byte streams between the connections.
func (p *Proxy) handleRequest() http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// ToDo: Implement request handling
	}

	return http.HandlerFunc(handler)
}
