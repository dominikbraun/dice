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
	"errors"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/store"
	"github.com/dominikbraun/dice/types"
)

var (
	ErrServiceNotFound      = errors.New("service could not be found")
	ErrServiceAlreadyExists = errors.New("a service with the given ID or name already exists")
)

// CreateService creates a new service with the provided name and stores
// the service in the key-value store. If the `Enable` option is set, the
// created service will be enabled immediately.
func (d *Dice) CreateService(name string, options types.ServiceCreateOptions) error {
	service, err := entity.NewService(name, options)
	if err != nil {
		return err
	}

	if ok, message := validateService(service); !ok {
		return errors.New(message)
	}

	isUnique, err := d.serviceIsUnique(service)

	if err != nil {
		return err
	} else if !isUnique {
		return ErrServiceAlreadyExists
	}

	if err := d.kvStore.CreateService(service); err != nil {
		return err
	}

	if options.Enable {
		return d.EnableService(entity.ServiceReference(service.ID))
	}

	return nil
}

// EnableService enables an existing service, making it available as request
// target. This function will update the service data and synchronize the
// service with the service registry.
func (d *Dice) EnableService(serviceRef entity.ServiceReference) error {
	service, err := d.findService(serviceRef)
	if err != nil {
		return err
	}

	if service == nil {
		return ErrServiceNotFound
	}

	service.IsEnabled = true

	if err := d.kvStore.UpdateService(service.ID, service); err != nil {
		return err
	}

	return d.synchronizeService(service, Enabling)
}

// DisableService disables a service, removing it as request target and
// therefore making it unavailable for any clients.
func (d *Dice) DisableService(serviceRef entity.ServiceReference) error {
	service, err := d.findService(serviceRef)
	if err != nil {
		return err
	}

	if service == nil {
		return ErrServiceNotFound
	}

	service.IsEnabled = false

	if err := d.kvStore.UpdateService(service.ID, service); err != nil {
		return err
	}

	return d.synchronizeService(service, Disabling)
}

// ServiceInfo returns user-relevant information for an existing service.
func (d *Dice) ServiceInfo(serviceRef entity.ServiceReference) (types.ServiceInfoOutput, error) {
	service, err := d.findService(serviceRef)
	if err != nil {
		return types.ServiceInfoOutput{}, err
	}

	if service == nil {
		return types.ServiceInfoOutput{}, ErrServiceNotFound
	}

	serviceInfo := types.ServiceInfoOutput{
		ID:              service.ID,
		Name:            service.Name,
		URLs:            service.URLs,
		TargetVersion:   service.TargetVersion,
		BalancingMethod: service.BalancingMethod,
		IsEnabled:       service.IsEnabled,
	}

	return serviceInfo, nil
}

// ListServices returns a list of stored services. By default, disabled
// services will be ignored. They only will be returned if the options say
// to do so.
func (d *Dice) ListServices(options types.ServiceListOptions) ([]types.ServiceInfoOutput, error) {
	filter := store.AllServicesFilter

	if !options.All {
		filter = func(service *entity.Service) bool {
			return service.IsEnabled
		}
	}

	services, err := d.kvStore.FindServices(filter)
	if err != nil {
		return nil, err
	}

	serviceList := make([]types.ServiceInfoOutput, len(services))

	for i, s := range services {
		info := types.ServiceInfoOutput{
			ID:              s.ID,
			Name:            s.Name,
			URLs:            s.URLs,
			TargetVersion:   s.TargetVersion,
			BalancingMethod: s.BalancingMethod,
			IsEnabled:       s.IsEnabled,
		}
		serviceList[i] = info
	}

	return serviceList, nil
}

// SetServiceURL sets or removes an URL from a given service. The update
// will be visible for the service registry and the Dice proxy instantly.
func (d *Dice) SetServiceURL(serviceRef entity.ServiceReference, url string, options types.ServiceURLOptions) error {
	service, err := d.findService(serviceRef)
	if err != nil {
		return err
	}

	if service == nil {
		return ErrServiceNotFound
	}

	if options.Delete {
		if err := service.RemoveURL(url); err != nil {
			return err
		}
	} else {
		if err := service.AddURL(url); err != nil {
			return err
		}
	}

	if err := d.kvStore.UpdateService(service.ID, service); err != nil {
		return err
	}

	return d.synchronizeService(service, Update)
}

// findService attempts to find a node in the key-value store that matches
// the reference. The ID has the highest priority, then the name is checked.
//
// If multiple services match, only the first one will be returned. If no
// services match, `nil` - and no error - will be returned.
func (d *Dice) findService(serviceRef entity.ServiceReference) (*entity.Service, error) {
	servicesByID, err := d.kvStore.FindServices(func(service *entity.Service) bool {
		return service.ID == string(serviceRef)
	})

	if err != nil {
		return nil, err
	} else if len(servicesByID) > 0 {
		return servicesByID[0], nil
	}

	servicesByName, err := d.kvStore.FindServices(func(service *entity.Service) bool {
		return service.Name == string(serviceRef)
	})

	if err != nil {
		return nil, err
	} else if len(servicesByName) > 0 {
		return servicesByName[0], nil
	}

	return nil, nil
}

// serviceIsUnique checks if a newly created service is unique. A service
// is unique if no service with equal identifiers has been found in the key
// value store.
func (d *Dice) serviceIsUnique(service *entity.Service) (bool, error) {
	storedService, err := d.findService(entity.ServiceReference(service.ID))

	if err != nil {
		return false, err
	} else if storedService != nil {
		return false, nil
	}

	if service.Name != "" {
		storedService, err = d.findService(entity.ServiceReference(service.Name))

		if err != nil {
			return false, err
		} else if storedService != nil {
			return false, nil
		}
	}

	return true, nil
}
