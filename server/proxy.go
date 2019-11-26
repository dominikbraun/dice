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

package server

import (
	"context"
	"github.com/dominikbraun/dice/registry"
	"net/http"
	"os"
)

type ProxyConfig struct {
	Address string `json:"address"`
	Logfile string `json:"logfile"`
}

type Proxy struct {
	config    ProxyConfig
	registry  *registry.ServiceRegistry
	server    *http.Server
	interrupt chan os.Signal
}

func NewProxy(config ProxyConfig, registry *registry.ServiceRegistry, quit chan os.Signal) *Proxy {
	p := Proxy{
		config:    config,
		registry:  registry,
		interrupt: quit,
	}

	p.server = &http.Server{
		Addr:    p.config.Address,
		Handler: nil,
	}

	return &p
}

func (p *Proxy) Run() chan<- error {
	errors := make(chan error)

	go func() {
		err := p.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errors <- err
		}
		close(errors)
	}()

	go func() {
		<-p.interrupt
		if err := p.server.Shutdown(context.Background()); err != nil {
			errors <- err
		}
	}()

	return errors
}

func (p *Proxy) handleRequest() http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// ToDo: Determine service and handle request
	}

	return http.HandlerFunc(handler)
}
