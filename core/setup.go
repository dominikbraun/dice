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
	"fmt"
	"github.com/dominikbraun/dice/api"
	"github.com/dominikbraun/dice/config"
	"github.com/dominikbraun/dice/controller"
	"github.com/dominikbraun/dice/healthcheck"
	"github.com/dominikbraun/dice/log"
	"github.com/dominikbraun/dice/proxy"
	"github.com/dominikbraun/dice/registry"
	"github.com/dominikbraun/dice/store"
	"os"
	"os/signal"
	"time"
)

// setupConfig parses the configuration file and sets all default values
// so that other components can rely on the keys. This step also powers
// Dice's zero-configuration ability.
func (d *Dice) setupConfig() error {
	var err error

	if d.config, err = config.NewConfig(configName); err != nil {
		return err
	}

	for key, value := range config.Defaults {
		d.config.SetDefault(key, value)
	}

	return nil
}

// setupKVStore opens or, if it doesn't exist, creates the key-value store.
func (d *Dice) setupKVStore() error {
	var err error

	if d.kvStore, err = store.NewKVStore(kvStorePath); err != nil {
		return err
	}

	return nil
}

// setupRegistry initializes the service registry. This is also the point
// where existing services and instances are acquainted to the registry.
func (d *Dice) setupRegistry() error {
	d.registry = registry.NewServiceRegistry()
	return nil
}

// setupHealthCheck initializes the default health checker. If no interval
// or timeout has been configured, Dice's default values will be used.
func (d *Dice) setupHealthCheck() error {
	var err error

	interval := d.config.GetInt("healthcheck-interval")
	timeout := d.config.GetInt("healthcheck-timeout")

	hcConfig := healthcheck.Config{
		Interval: time.Duration(interval) * time.Millisecond,
		Timeout:  time.Duration(timeout) * time.Millisecond,
	}

	if d.healthCheck, err = healthcheck.New(hcConfig, &d.registry.Services); err != nil {
		return err
	}

	return nil
}

// setupController creates a new Controller instance that utilizes Dice
// itself as a controller target. It will be used by the API server.
func (d *Dice) setupController() error {
	d.controller = controller.New(d)
	return nil
}

// setupAPIServer configures the API server, however it won't be started.
func (d *Dice) setupAPIServer() error {
	port := d.config.GetString("api-server-port")
	address := fmt.Sprintf(":%v", port)

	logfile := d.config.GetString("api-server-logfile")

	serverConfig := api.ServerConfig{
		Address: address,
		Logfile: logfile,
	}

	d.apiServer = api.NewServer(serverConfig, d.controller)

	return nil
}

// setupProxy configures the proxy server, which won't be started either.
func (d *Dice) setupProxy() error {
	port := d.config.GetString("proxy-port")
	address := fmt.Sprintf(":%v", port)

	logfile := d.config.GetString("proxy-logfile")

	proxyConfig := proxy.Config{
		Address: address,
		Logfile: logfile,
	}

	d.proxy = proxy.New(proxyConfig, d.registry)

	return nil
}

// setupLogger sets up the logger as well as the logfile it will be using.
func (d *Dice) setupLogger() error {
	logfile := d.config.GetString("dice-logfile")

	file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	d.logger = log.NewLogger(file, log.InfoLevel)

	return nil
}

// setupInterrupt creates the interrupt channel. It will be notified if a
// system signal (SIGINT) is sent to the Dice executable.
func (d *Dice) setupInterrupt() error {
	d.interrupt = make(chan os.Signal)
	signal.Notify(d.interrupt, os.Interrupt)

	return nil
}
