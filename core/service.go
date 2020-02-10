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
	"github.com/dominikbraun/dice/registry"
	"github.com/dominikbraun/dice/store"
	"github.com/dominikbraun/dice/types"
	"strings"
)

var (
	ErrServiceNotFound      = errors.New("service could not be found")
	ErrServiceAlreadyExists = errors.New("a service with the given ID or name already exists")
	ErrServiceURLExists     = errors.New("one or more of the specified URLs already exists")
)

// CreateService creates a new service with the provided name and stores
// the service in the key-value store. If the `Enable` option is set, the
// created service will be enabled immediately.
func (d *Dice) CreateService(name string, options types.ServiceCreateOptions) error {
	service, err := entity.NewService(name, options)
	if err != nil {
		return err
	}

	ok, err := d.urlsAreValid(service)
	if err != nil {
		return err
	}

	if !ok {
		return ErrServiceURLExists
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

	if err := d.registry.Register(service, d.buildRegistryService); err != nil {
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
	} else if service == nil {
		return ErrServiceNotFound
	}

	service.IsEnabled = true

	if err := d.kvStore.UpdateService(service.ID, service); err != nil {
		return err
	}

	return d.registry.Update(func(s *registry.Service) error {
		if s.Entity.ID == service.ID {
			s.Entity.IsEnabled = true
		}
		return nil
	})
}

// DisableService disables a service, removing it as request target and
// therefore making it unavailable for any clients.
func (d *Dice) DisableService(serviceRef entity.ServiceReference) error {
	service, err := d.findService(serviceRef)

	if err != nil {
		return err
	} else if service == nil {
		return ErrServiceNotFound
	}

	service.IsEnabled = false

	if err := d.kvStore.UpdateService(service.ID, service); err != nil {
		return err
	}

	return d.registry.Update(func(s *registry.Service) error {
		if s.Entity.ID == service.ID {
			s.Entity.IsEnabled = false
		}
		return nil
	})
}

// UpdateService updates a service whose instances have already been deployed
// under specific version tags. That is, all instances whose versions do not
// match the targetVersion will be detached. Instances that have a matching
// version will be attached.
func (d *Dice) UpdateService(serviceRef entity.ServiceReference, targetVersion string) error {
	service, err := d.findService(serviceRef)

	if err != nil {
		return err
	} else if service == nil {
		return ErrServiceNotFound
	}

	attachableInstances, err := d.kvStore.FindInstances(func(instance *entity.Instance) bool {
		return instance.Version == strings.Trim(targetVersion, " ")
	})

	if err != nil {
		return err
	}

	for _, i := range attachableInstances {
		// AttachInstance and DetachInstance will search the KV store entry
		// again in order to create an instance, change it and write it back.
		// ToDo: Avoid loading instances from the KV store twice.
		if err := d.AttachInstance(entity.InstanceReference(i.ID)); err != nil {
			return err
		}
	}

	detachableInstances, err := d.kvStore.FindInstances(func(instance *entity.Instance) bool {
		return instance.Version != strings.Trim(targetVersion, " ")
	})

	if err != nil {
		return err
	}

	for _, i := range detachableInstances {
		if err := d.DetachInstance(entity.InstanceReference(i.ID)); err != nil {
			return err
		}
	}

	return nil
}

// ServiceInfo returns user-relevant information for an existing service.
func (d *Dice) ServiceInfo(serviceRef entity.ServiceReference) (types.ServiceInfoOutput, error) {
	service, err := d.findService(serviceRef)

	if err != nil {
		return types.ServiceInfoOutput{}, err
	} else if service == nil {
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
	} else if service == nil {
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

	if options.Delete {
		if err := d.registry.UnregisterServiceURL(url); err != nil {
			return err
		}
	} else {
		if err := d.registry.RegisterServiceURL(service.ID, url); err != nil {
			return err
		}
	}

	return d.registry.Update(func(s *registry.Service) error {
		if s.Entity.ID == service.ID {
			s.Entity.URLs = service.URLs
		}
		return nil
	})
}

// urlsAreValid indicates whether a services' URLs are valid and unique
// so that it can be used safely. This check should be performed before
// the service entity gets persisted.
func (d *Dice) urlsAreValid(service *entity.Service) (bool, error) {
	servicesByURL, err := d.kvStore.FindServices(func(s *entity.Service) bool {
		for _, u := range s.URLs {
			for _, su := range service.URLs {
				if u == su {
					return true
				}
			}
		}
		return false
	})

	if err != nil {
		return false, err
	}
	isValid := len(servicesByURL) == 0

	return isValid, nil
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
