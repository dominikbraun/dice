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

// Package registry provides a global registry for fast access on services.
// Unlike the storage package, the registry package provides data in a
// representation designed for working with the data.
package registry

import (
	"fmt"

	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/scheduler"
)

// Service is the registry representation of a service. It does not only
// include the entity itself but also a scheduler and its deployments.
type Service struct {
	entity      *entity.Service
	scheduler   scheduler.Scheduler
	deployments []entity.Deployment
}

// ServiceRegistry is a registry that provides quick access to services,
// their responsible schedulers and their current deployments.
type ServiceRegistry struct {
	services map[string]Service
	hosts    map[string]string
}

// NewServiceRegistry creates a new ServiceRegistry instance.
func NewServiceRegistry() *ServiceRegistry {
	sr := ServiceRegistry{
		services: make(map[string]Service),
		hosts:    make(map[string]string),
	}

	return &sr
}

// Setup builds the service registry by iterating over a list of service
// entities and invoking a function for each entity. This function should
// transform a given entity.Service into a complete registry.Service.
func (sr *ServiceRegistry) Setup(services []*entity.Service, transform func(*entity.Service) Service) error {
	for _, s := range services {
		service := transform(s)

		if err := sr.AddService(service); err != nil {
			return err
		}
	}

	return nil
}

// AddService adds a new service. Returns an error if it already exists.
func (sr *ServiceRegistry) AddService(service Service) error {
	serviceID := service.entity.ID

	if _, exists := sr.services[serviceID]; exists {
		return fmt.Errorf("service %v is already registered", serviceID)
	}
	sr.services[serviceID] = service

	for _, h := range service.entity.Config.Hosts {
		sr.hosts[h] = service.entity.ID
	}

	return nil
}

// UpdateService updates an existing service entry. Returns an error if
// the service is not registered.
func (sr *ServiceRegistry) UpdateService(serviceID string, service Service) error {
	if _, exists := sr.services[serviceID]; !exists {
		return fmt.Errorf("service %v is not registered", serviceID)
	}
	sr.services[serviceID] = service

	return nil
}

// RemoveService removes a service and all host entries from the registry.
func (sr *ServiceRegistry) RemoveService(serviceID string) error {
	if _, exists := sr.services[serviceID]; !exists {
		return fmt.Errorf("service %v is not registered", serviceID)
	}

	for host, id := range sr.hosts {
		if id == serviceID {
			delete(sr.hosts, host)
		}
	}
	delete(sr.services, serviceID)

	return nil
}

// ServiceExits determines if a service is registered.
func (sr *ServiceRegistry) ServiceExists(serviceID string) bool {
	_, exists := sr.services[serviceID]
	return exists
}

// LookupService searches a service using a given service ID.
func (sr *ServiceRegistry) LookupService(serviceID string) (Service, error) {
	if _, exists := sr.services[serviceID]; !exists {
		return Service{}, fmt.Errorf("service %v is not registered", serviceID)
	}

	return sr.services[serviceID], nil
}

// LookupServiceByHost searches a service by one of its host addresses.
func (sr *ServiceRegistry) LookupServiceByHost(host string) (Service, error) {
	serviceID, exists := sr.hosts[host]

	if !exists {
		return Service{}, fmt.Errorf("no service with host %v registered", host)
	}

	return sr.LookupService(serviceID)
}

// AddServiceHost adds a new host to the registry. Note that this function
// does not add the host to the service entity.
func (sr *ServiceRegistry) AddServiceHost(host, serviceID string) error {
	if _, exists := sr.hosts[host]; exists {
		return fmt.Errorf("host %v is already registered", host)
	}
	sr.hosts[host] = serviceID

	return nil
}

// RemoveServiceHost removes an existing host entry from the registry. Note
// that this function does not remove the host from the service entity.
func (sr *ServiceRegistry) RemoveServiceHost(host string) error {
	if _, exists := sr.hosts[host]; !exists {
		return fmt.Errorf("host %v is not registered", host)
	}
	delete(sr.hosts, host)

	return nil
}
