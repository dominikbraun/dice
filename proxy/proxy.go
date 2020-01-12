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

// Package proxy provides a reverse proxy. Its job is to accept incoming
// requests, find a service instance and forward the request to it.
package proxy

import (
	"context"
	"fmt"
	"github.com/dominikbraun/dice/registry"
	"net/http"
)

// Config concludes properties that are configurable by the user.
type Config struct {
	Address string `json:"address"`
	Logfile string `json:"logfile"`
}

// Proxy is a reverse proxy that accepts incoming requests for all services,
// looks up the responsible service in the registry and proxies the request
// for to an instance of that service.
//
// Proxy only uses read-only access on ServiceRegistry.
type Proxy struct {
	config   Config
	registry *registry.ServiceRegistry
	server   *http.Server
}

// New creates a new Proxy instance and sets up a ready-to-go HTTP server.
func New(config Config, registry *registry.ServiceRegistry) *Proxy {
	p := Proxy{
		config:   config,
		registry: registry,
	}

	p.server = &http.Server{
		Addr:    p.config.Address,
		Handler: p.handleRequest(),
	}

	return &p
}

// Run starts the proxy, accepting incoming requests on the configured port.
func (p *Proxy) Run() error {
	err := p.server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown attempts a graceful shutdown of the proxy server. It will wait
// for all open connections to finish and stops the proxy subsequently.
func (p *Proxy) Shutdown() error {
	err := p.server.Shutdown(context.Background())
	_ = p.server.Close()

	return err
}

// handleRequest processes an incoming request. After looking up the desired
// service in the service registry, the provided scheduler will be used to
// obtain a service instance. Proxy will then establish a connection to that
// instance, forward the request to it and send the response back to the client.
func (p *Proxy) handleRequest() http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		service, ok := p.registry.LookupService(r.Host)
		_, _ = w.Write([]byte(fmt.Sprintf("Service: %#+v\n", service)))
		_, _ = w.Write([]byte(fmt.Sprintf("Is registered: %v\n", ok)))

		if !ok || !service.Entity.IsEnabled || service.Scheduler == nil {
			p.displayError(w, r, http.StatusServiceUnavailable, "Service Unavailable")
			return
		}

		_, err := service.Scheduler.Next()
		if err != nil {
			p.displayError(w, r, http.StatusServiceUnavailable, "Service Unavailable")
			return
		}

		// ToDo: Forward request and return response.
	}

	return http.HandlerFunc(handler)
}

// displayError returns an error response to the client by setting the provided
// HTTP status code and displaying the desired message.
func (p *Proxy) displayError(w http.ResponseWriter, r *http.Request, status int, message string) {
	const template = `
<body style="text-align: center">
	<h1 style="font-family: arial">Error %d: %s</h1>
	<hr />
	<p style="font-family: arial">Dice</p>`

	body := fmt.Sprintf(template, status, message)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
