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

// Package registry provides the service registry and the route registry.
//
// While the core package as well as the store package represent the data
// statically and storage-oriented, the registries provide a representation
// required at runtime: In-memory, dynamic and quickly accessible.
package registry

import (
	"errors"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/log"
)

type (
	NodeUpdater     func(node *entity.Node) error
	ServiceUpdater  func(service *entity.Service) error
	InstanceUpdater func(instance *entity.Instance) error
)

var (
	ErrUnregisteredService      = errors.New("service is not registered")
	ErrServiceAlreadyRegistered = errors.New("service is already registered")
	ErrServiceNotRemovable      = errors.New("service has attached instances on an attached node")
	ErrUnregisteredDeployment   = errors.New("deployment is not registered")
	ErrDeploymentNotRemovable   = errors.New("deployed instance is attached on an attached node")
)

// ServiceRegistry is the global registry for all services known to Dice.
// Its purpose is to provide quick access to deployment information about
// a particular service at runtime.
//
// When the proxy asks for the service associated with a particular host,
// the ServiceRegistry looks up that service and returns all information
// like deployments and the scheduler (find more in the `Service` docs).
//
// ServiceRegistry also offers methods for updating existing service data
// and for registering new services or service deployments at runtime.
type ServiceRegistry struct {
	Services      map[string]*Service
	routeRegistry RouteRegistry
	logger        log.Logger
}

// NewServiceRegistry creates a new ServiceRegistry instance that writes
// to a given log.Logger. The new instance has to be initialized with all
// stored services on startup, see `Register`.
func NewServiceRegistry(logger log.Logger) *ServiceRegistry {
	sr := ServiceRegistry{
		Services:      make(map[string]*Service),
		routeRegistry: NewRouteRegistry(),
		logger:        logger,
	}

	return &sr
}

// Register registers a new service. The build function should return a
// fully initialized registry.Service instance, including deployments and
// scheduler.
func (sr *ServiceRegistry) Register(entity *entity.Service, build func(*entity.Service) (*Service, error)) error {
	service, err := build(entity)
	if err != nil {
		return err
	}

	return sr.RegisterService(service, false)
}

// RegisterService registers a new service. Returns an error if the service
// is already registered, unless force is set to `true`.
func (sr *ServiceRegistry) RegisterService(service *Service, force bool) error {
	serviceID := service.Entity.ID

	if _, exists := sr.Services[serviceID]; exists {
		if !force {
			return ErrServiceAlreadyRegistered
		}
	}

	for _, r := range service.Entity.URLs {
		if err := sr.routeRegistry.RegisterRoute(r, serviceID, force); err != nil {
			return err
		}
	}

	sr.Services[serviceID] = service
	return nil
}

// UnregisterService removes a registered service from the registry. Returns
// an error if the service has attached instances on attached nodes, unless
// force is set to `true`.
func (sr *ServiceRegistry) UnregisterService(serviceID string, force bool) error {
	if _, exists := sr.Services[serviceID]; !exists {
		return ErrUnregisteredService
	}

	if !force {
		for _, d := range sr.Services[serviceID].Deployments {
			if !d.isRemovable() {
				return ErrServiceNotRemovable
			}
		}
	}

	for _, r := range sr.Services[serviceID].Entity.URLs {
		if err := sr.routeRegistry.UnregisterRoute(r); err != nil {
			return err
		}
	}

	delete(sr.Services, serviceID)
	return nil
}

// LookupService looks up the service available under a given route. The
// second return value indicates whether the service could be found or not.
func (sr *ServiceRegistry) LookupService(host string) (*Service, bool) {
	serviceID, exists := sr.routeRegistry.LookupServiceID(host)
	if !exists {
		return &Service{}, false
	}

	if service, exists := sr.Services[serviceID]; exists {
		return service, true
	}
	sr.logger.Warnf("service %s registered in router but not in registry", serviceID)

	return &Service{}, false
}

// Update is the public API for accessing the registry services and applying
// an update function on each of them. This function may be used to update the
// service entity itself or some node or instance information.
//
// Update should be the only way for other components to gain write-access to
// the registry's internal services.
func (sr *ServiceRegistry) Update(updateFunc func(service *Service) error) error {
	for _, s := range sr.Services {
		if err := updateFunc(s); err != nil {
			return err
		}
	}

	return nil
}

// RegisterServiceURL registers a new public URL for a service. Returns an
// error of the given URL already exists for this or another service.
func (sr *ServiceRegistry) RegisterServiceURL(serviceID, url string) error {
	return sr.routeRegistry.RegisterRoute(url, serviceID, false)
}

// UnregisterServiceURL removes a public URL from the registry. Unregistering
// an URL will cause Dice to return an error for requests related to that URL.
func (sr *ServiceRegistry) UnregisterServiceURL(url string) error {
	return sr.routeRegistry.UnregisterRoute(url)
}

// RegisterDeployment registers new service deployment. Returns an error
// if the stored service in the `Instance` field is not registered yet.
func (sr *ServiceRegistry) RegisterDeployment(deployment Deployment) error {
	serviceID := deployment.Instance.ServiceID

	if _, exists := sr.Services[serviceID]; !exists {
		return ErrUnregisteredService
	}

	service := sr.Services[serviceID]
	service.Deployments = append(service.Deployments, deployment)

	service.Scheduler.UpdateDeployments(service.Deployments)

	return nil
}

// UnregisterDeployment removes a deployment from the registry. Returns
// an error if the instance is attach on an attached node, unless force
// is set to `true`.
func (sr *ServiceRegistry) UnregisterDeployment(deployment Deployment, force bool) error {
	serviceID := deployment.Instance.ServiceID

	if _, exists := sr.Services[serviceID]; !exists {
		return ErrUnregisteredService
	}

	index, err := sr.indexOfDeployment(serviceID, deployment)
	if err != nil {
		return err
	} else if index == -1 {
		return ErrUnregisteredDeployment
	}

	if !force {
		if !deployment.isRemovable() {
			return ErrDeploymentNotRemovable
		}
	}

	service := sr.Services[serviceID]
	service.Deployments[index] = service.Deployments[len(service.Deployments)-1]
	service.Deployments = service.Deployments[:len(service.Deployments)-1]

	service.Scheduler.UpdateDeployments(service.Deployments)

	return nil
}

// indexOfDeployment determines the index of a service's deployment. The
// given deployment is considered equal to another deployment if its node
// ID and instance ID are the same. Returns -1 if no deployment matches.
func (sr *ServiceRegistry) indexOfDeployment(serviceID string, deployment Deployment) (int, error) {
	if _, exists := sr.Services[serviceID]; !exists {
		return 0, ErrUnregisteredService
	}

	for i, d := range sr.Services[serviceID].Deployments {
		if d.equals(deployment) {
			return i, nil
		}
	}

	return -1, nil
}
