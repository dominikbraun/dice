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
	"github.com/dominikbraun/dice/log"
	"github.com/dominikbraun/dice/proxy"
	"github.com/dominikbraun/dice/registry"
	"github.com/dominikbraun/dice/store"
	"os"
	"os/signal"
)

const (
	configName  string = "dice"
	kvStorePath string = "./dice-store"
	logfilePath string = "./dice.log"
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
	config    config.Reader
	kvStore   store.EntityStore
	registry  *registry.ServiceRegistry
	interrupt chan os.Signal
	apiServer *api.Server
	proxy     *proxy.Proxy
	logger    log.Logger
}

// NewDice creates a new Dice instances and initializes all components.
func NewDice() (*Dice, chan<- os.Signal, error) {
	var d Dice
	var err error

	if d.config, err = config.NewConfig(configName); err != nil {
		return nil, nil, err
	}

	if d.kvStore, err = store.NewKV(kvStorePath); err != nil {
		return nil, nil, err
	}

	d.registry = registry.NewServiceRegistry()
	d.interrupt = make(chan os.Signal)

	d.apiServer = api.NewServer(api.ServerConfig{})
	d.proxy = proxy.New(proxy.Config{}, d.registry)

	logfile, err := os.OpenFile(logfilePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, nil, err
	}

	d.logger = log.NewLogger(logfile, log.InfoLevel)

	signal.Notify(d.interrupt, os.Interrupt)

	return &d, d.interrupt, nil
}

// Run starts the API and proxy servers. To shut them down gracefully, send
// a signal through the interrupt channel returned by NewDice.
func (d *Dice) Run() error {
	errors := make(chan error)

	go func() {
		if err := d.proxy.Run(d.interrupt); err != nil {
			errors <- err
		}
	}()

	go func() {
		if err := d.apiServer.Run(d.interrupt); err != nil {
			errors <- err
		}
	}()

	err := <-errors
	return err
}
