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
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/healthcheck"
	"github.com/dominikbraun/dice/log"
	"github.com/dominikbraun/dice/proxy"
	"github.com/dominikbraun/dice/registry"
	"github.com/dominikbraun/dice/scheduler"
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
	config       config.Reader
	reloadConfig chan bool
	logger       log.Logger
	kvStore      store.EntityStore
	registry     *registry.ServiceRegistry
	healthCheck  *healthcheck.HealthCheck
	controller   *controller.Controller
	interrupt    chan os.Signal
	apiServer    *api.Server
	proxy        *proxy.Proxy
}

// NewDice creates a new Dice instance and sets up all components.
func NewDice() (*Dice, error) {
	var d Dice

	if err := d.setup(); err != nil {
		return nil, err
	}

	return &d, nil
}

// setup runs the Dice setup by invoking all setup* methods.
func (d *Dice) setup() error {
	steps := []func() error{
		d.setupConfig,
		d.setupReloadConfig,
		d.setupLogger,
		d.setupKVStore,
		d.setupRegistry,
		d.setupHealthCheck,
		d.setupController,
		d.setupAPIServer,
		d.setupProxy,
		d.setupInterrupt,
	}

	for _, setup := range steps {
		if err := setup(); err != nil {
			return err
		}
	}

	return nil
}

// Run starts the API and proxy servers. To shut them down gracefully, send
// an interrupt signal (SIGINT) to the Dice executable. If an error happens
// while running one of the servers, Dice will be stopped entirely.
func (d *Dice) Run() error {
	if err := d.initializeRegistry(); err != nil {
		return err
	}

	for {
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
			return nil

		case reload := <-d.reloadConfig:
			d.logger.Info("reloading Dice")

			if reload {
				if err := d.proxy.Shutdown(); err != nil {
					d.logger.Errorf("proxy shutdown error: %v", err)
				}
				if err := d.apiServer.Shutdown(); err != nil {
					d.logger.Errorf("API server shutdown error: %v", err)
				}
				if err := d.setup(); err != nil {
					return err
				}
			}

		case err := <-errors:
			return err
		}
	}
}

// initializeServices initializes all services and makes them available for
// load balancing. This is done by populating the service registry with all
// services, their deployments and the responsible scheduler.
//
// This method only sets up the services for the registry at startup time. At
// runtime, services and deployments will be registered by core methods like
// CreateService using the exact same mechanisms.
//
// ToDo: Clarify how errors during initialization should be handled.
func (d *Dice) initializeRegistry() error {
	services, err := d.kvStore.FindServices(store.AllServicesFilter)
	if err != nil {
		return err
	}

	for _, s := range services {
		registryService, err := d.buildRegistryService(s)
		if err != nil {
			return err
		}

		if err := d.registry.RegisterService(registryService, false); err != nil {
			if err != registry.ErrRouteAlreadyRegistered {
				return err
			}
		}
	}

	return nil
}

// buildRegistryService takes a service entity and creates a registry.Service
// instance by searching the instances and the nodes they've been deployed to.
//
// The created registry.Service includes information about deployed instances
// of the particular service and provides a scheduler as well.
//
// See the registry.Service docs for further explanations.
func (d *Dice) buildRegistryService(service *entity.Service) (*registry.Service, error) {
	registryService := registry.Service{
		Entity: service,
	}

	instances, err := d.kvStore.FindInstances(func(i *entity.Instance) bool {
		return i.ServiceID == service.ID
	})
	if err != nil {
		return &registryService, err
	}

	registryService.Deployments = make([]registry.Deployment, len(instances))

	for i, inst := range instances {
		node, err := d.kvStore.FindNode(inst.NodeID)
		if err != nil {
			return &registryService, err
		}

		registryService.Deployments[i] = registry.Deployment{
			Node:     node,
			Instance: inst,
		}
	}

	serviceScheduler, err := scheduler.New(registryService.Deployments, scheduler.BalancingMethod(service.BalancingMethod))
	if err != nil {
		return &registryService, err
	}

	registryService.Scheduler = serviceScheduler
	return &registryService, nil
}
