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
	"fmt"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/scheduler"
	"github.com/dominikbraun/dice/storage"
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
	config ProxyConfig
	memory storage.Entity
	server *http.Server

	// hostRegistry is a mapping of a hostname against a service ID.
	hostRegistry map[string]string

	// serviceRegistry is a mapping of a service ID against a service.
	serviceRegistry map[string]Service
}

// Service represents an entity.Service from the proxy's point of view. This
// means that it holds the original service data, an associated scheduler and
// most importantly a list of deployed service instances.
//
// The data about the services - such as existing services and their instances -
// will be read from Dice's memory storage and then stored in the service map.
type Service struct {
	entity      *entity.Service
	scheduler   scheduler.Scheduler
	deployments []entity.Deployment
}

// NewProxy creates a new Proxy instance and returns a reference to it.
func NewProxy(config ProxyConfig, memory storage.Entity) *Proxy {
	p := Proxy{
		config:          config,
		memory:          memory,
		hostRegistry:    make(map[string]string),
		serviceRegistry: make(map[string]Service),
	}

	p.server = &http.Server{
		Addr:    p.config.Address,
		Handler: p.handleRequest(),
	}

	return &p
}

// initRegistries initializes the host- and service registries. It reads
// all stored services from the in-memory storage and populates the maps.
func (p *Proxy) initRegistries() error {
	services, err := p.memory.FindAll(storage.Service)
	if err != nil {
		return err
	}

	for _, s := range services {
		e := s.(*entity.Service)

		for _, h := range e.Config.Hosts {
			p.hostRegistry[h] = e.ID
		}

		service := Service{
			entity:      e,
			deployments: nil,
		}

		method := scheduler.BalancingMethod(e.Config.BalancingMethod)
		service.scheduler, err = scheduler.NewScheduler(method, &service.deployments)
		if err != nil {
			return err
		}

		p.serviceRegistry[service.entity.ID] = service
	}

	return nil
}

// lookupService searches an entry for the given host in the host registry. If
// it exists, it searches and returns the service with the associated ID.
func (p *Proxy) lookupService(host string) (Service, error) {
	serviceID, exists := p.hostRegistry[host]
	if !exists {
		return Service{}, fmt.Errorf("host entry for %v not found", host)
	}

	service, exists := p.serviceRegistry[serviceID]
	if !exists {
		return Service{}, fmt.Errorf("service for %v not found", host)
	}

	return service, nil
}

// Run starts the proxy server. It will listen to the specified port and handle
// incoming requests, sending errors through the returned channel.
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

// handleRequest implements the core request handling logic which includes the
// transfer/copy of the byte streams between the connections.
func (p *Proxy) handleRequest() http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// ToDo: Implement request handling
	}

	return http.HandlerFunc(handler)
}
