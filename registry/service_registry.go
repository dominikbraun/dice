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

package registry

import (
	"errors"
	"github.com/dominikbraun/dice/entity"
)

type UnregisterMode uint

const (
	SoftUnregister UnregisterMode = 0
	HardUnregister UnregisterMode = 1
)

type (
	NodeFilter       func(node *entity.Node) bool
	ServiceFilter    func(service *entity.Service) bool
	InstanceFilter   func(instance *entity.Instance) bool
	DeploymentFilter func(deployment Deployment) bool
)

type (
	NodeUpdater       func(node *entity.Node) error
	ServiceUpdater    func(service *entity.Service) error
	InstanceUpdater   func(instance *entity.Instance) error
	DeploymentUpdater func(deployment Deployment) error
)

var (
	ErrUnregisteredService       = errors.New("service is not registered")
	ErrServiceAlreadyRegistered  = errors.New("service is already registered")
	ErrServiceNotRemovable       = errors.New("service has attached instances on an attached node")
	ErrUnregisteredDeployment    = errors.New("deployment is not registered")
	ErrDeploymentNotRemovable    = errors.New("deployed instance is attached on an attached node")
	ErrUnregisteredHostname      = errors.New("hostname is not registered")
	ErrHostnameAlreadyRegistered = errors.New("hostname is already registered")
	ErrInvalidClosure            = errors.New("the provided closure is invalid")
)

type ServiceRegistry struct {
	Services  map[string]Service
	Hostnames map[string]string
}

func NewServiceRegistry() *ServiceRegistry {
	sr := ServiceRegistry{
		Services:  make(map[string]Service),
		Hostnames: make(map[string]string),
	}

	return &sr
}

func (sr *ServiceRegistry) Register(entity *entity.Service, build func(*entity.Service) Service) error {
	service := build(entity)
	return sr.RegisterService(service)
}

func (sr *ServiceRegistry) UpdateNodes(filter NodeFilter, updater NodeUpdater) error {
	if filter == nil || updater == nil {
		return ErrInvalidClosure
	}

	for _, s := range sr.Services {
		for _, d := range s.Deployments {
			if filter(d.Node) {
				if err := updater(d.Node); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (sr *ServiceRegistry) UpdateServices(filter ServiceFilter, updater ServiceUpdater) error {
	if filter == nil || updater == nil {
		return ErrInvalidClosure
	}

	for _, s := range sr.Services {
		if filter(s.Entity) {
			if err := updater(s.Entity); err != nil {
				return err
			}
		}
	}

	return nil
}

func (sr *ServiceRegistry) UpdateInstances(filter InstanceFilter, updater InstanceUpdater) error {
	if filter == nil || updater == nil {
		return ErrInvalidClosure
	}

	for _, s := range sr.Services {
		for _, d := range s.Deployments {
			if filter(d.Instance) {
				if err := updater(d.Instance); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (sr *ServiceRegistry) UpdateDeployments(filter DeploymentFilter, updater DeploymentUpdater) error {
	if filter == nil || updater == nil {
		return ErrInvalidClosure
	}

	for _, s := range sr.Services {
		for _, d := range s.Deployments {
			if filter(d) {
				if err := updater(d); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (sr *ServiceRegistry) RegisterService(service Service) error {
	serviceID := service.Entity.ID

	if _, exists := sr.Services[serviceID]; exists {
		return ErrServiceAlreadyRegistered
	}

	for _, h := range service.Entity.Hostnames {
		if err := sr.RegisterHostname(h, serviceID); err != nil {
			return err
		}
	}

	sr.Services[serviceID] = service
	return nil
}

func (sr *ServiceRegistry) UnregisterService(serviceID string, mode UnregisterMode) error {
	if _, exists := sr.Services[serviceID]; !exists {
		return ErrUnregisteredService
	}

	if mode != HardUnregister {
		for _, d := range sr.Services[serviceID].Deployments {
			if !d.isRemovable() {
				return ErrServiceNotRemovable
			}
		}
	}

	for _, h := range sr.Services[serviceID].Entity.Hostnames {
		if err := sr.UnregisterHostname(h); err != nil {
			return err
		}
	}

	delete(sr.Services, serviceID)
	return nil
}

func (sr *ServiceRegistry) RegisterDeployment(deployment Deployment) error {
	serviceID := deployment.Instance.ServiceID

	if _, exists := sr.Services[serviceID]; !exists {
		return ErrUnregisteredService
	}

	deployments := sr.Services[serviceID].Deployments
	deployments = append(deployments, deployment)

	return nil
}

func (sr *ServiceRegistry) UnregisterDeployment(deployment Deployment, mode UnregisterMode) error {
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

	if mode != HardUnregister {
		if !deployment.isRemovable() {
			return ErrDeploymentNotRemovable
		}
	}

	deployments := sr.Services[serviceID].Deployments
	deployments[index] = deployments[len(deployments)-1]
	deployments = deployments[:len(deployments)-1]

	return nil
}

func (sr *ServiceRegistry) RegisterHostname(hostname string, serviceID string) error {
	if _, exists := sr.Hostnames[hostname]; exists {
		return ErrHostnameAlreadyRegistered
	}

	if _, exists := sr.Services[serviceID]; !exists {
		return ErrUnregisteredService
	}
	sr.Hostnames[hostname] = serviceID

	return nil
}

func (sr *ServiceRegistry) UnregisterHostname(hostname string) error {
	if _, exists := sr.Hostnames[hostname]; !exists {
		return ErrUnregisteredHostname
	}
	delete(sr.Hostnames, hostname)

	return nil
}

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
