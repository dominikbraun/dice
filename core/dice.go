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

// Package core provides the Dice load balancer and its methods.
package core

import (
	"github.com/dominikbraun/dice/api"
	"github.com/dominikbraun/dice/config"
	"github.com/dominikbraun/dice/controller"
	"github.com/dominikbraun/dice/healthcheck"
	"github.com/dominikbraun/dice/log"
	"github.com/dominikbraun/dice/proxy"
	"github.com/dominikbraun/dice/registry"
	"github.com/dominikbraun/dice/store"
	"os"
)

const (
	configName string = "dice"
)

// Dice represents the Dice load balancer and wires up all the components.
//
// Most importantly, this type consists of:
// - a key-value store for simply persisting domain entities
// - a registry that manages all services and their instances
// - an API server that exposes a REST API for managing Dice
// - a proxy server that will receive and balance all requests
// - some utility types for configuration parsing, logging etc.
//
// Some deeper explanations can be found at the corresponding components.
type Dice struct {
	config      config.Reader
	kvStore     store.EntityStore
	registry    *registry.ServiceRegistry
	healthCheck *healthcheck.HealthCheck
	controller  *controller.Controller
	interrupt   chan os.Signal
	apiServer   *api.Server
	proxy       *proxy.Proxy
	logger      log.Logger
}

// NewDice creates a new Dice instance and invokes all setup methods.
func NewDice() (*Dice, error) {
	var d Dice

	steps := []func() error{
		d.setupConfig,
		d.setupKVStore,
		d.setupRegistry,
		d.setupHealthCheck,
		d.setupController,
		d.setupAPIServer,
		d.setupProxy,
		d.setupLogger,
		d.setupInterrupt,
	}

	for _, setup := range steps {
		if err := setup(); err != nil {
			return nil, err
		}
	}

	return &d, nil
}

// Run starts the API and proxy servers. To shut them down gracefully, send
// an interrupt signal (SIGINT) to the Dice executable. If an error happens
// while running one of the servers, Dice will be stopped entirely.
func (d *Dice) Run() error {
	errors := make(chan error)

	go func() {
		if err := d.proxy.Run(); err != nil {
			errors <- err
		}
	}()

	go func() {
		if err := d.apiServer.Run(); err != nil {
			errors <- err
		}
	}()

	select {
	case <-d.interrupt:
		if err := d.proxy.Shutdown(); err != nil {
			d.logger.Errorf("Proxy shutdown error: %v", err)
		}
		if err := d.apiServer.Shutdown(); err != nil {
			d.logger.Errorf("API server shutdown error: %v", err)
		}
	case err := <-errors:
		return err
	}

	return nil
}
