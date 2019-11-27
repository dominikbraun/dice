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

package core

import (
	"errors"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/types"
)

type ServiceReference string

var (
	ErrServiceNotFound      = errors.New("service could not be found")
	ErrServiceAlreadyExists = errors.New("a service with the given ID or name already exists")
)

func (d *Dice) ServiceCreate(name string, options types.ServiceCreateOptions) error {
	service, err := entity.NewService(name, options)
	if err != nil {
		return err
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
		return d.ServiceEnable(ServiceReference(service.ID))
	}

	return nil
}

func (d *Dice) ServiceEnable(serviceRef ServiceReference) error {
	service, err := d.findService(serviceRef)
	if err != nil {
		return err
	}

	if service == nil {
		return ErrServiceNotFound
	}

	service.IsEnabled = true

	return d.kvStore.UpdateService(service.ID, service)
}

func (d *Dice) ServiceDisable(serviceRef ServiceReference) error {
	service, err := d.findService(serviceRef)
	if err != nil {
		return err
	}

	if service == nil {
		return ErrServiceNotFound
	}

	service.IsEnabled = false

	return d.kvStore.UpdateService(service.ID, service)
}

func (d *Dice) ServiceInfo(serviceRef ServiceReference) (types.ServiceInfoOutput, error) {
	var serviceInfo types.ServiceInfoOutput

	service, err := d.findService(serviceRef)
	if err != nil {
		return serviceInfo, err
	}

	serviceInfo.Populate(service)

	return serviceInfo, nil
}

func (d *Dice) findService(serviceRef ServiceReference) (*entity.Service, error) {
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

func (d *Dice) serviceIsUnique(service *entity.Service) (bool, error) {
	service, err := d.findService(ServiceReference(service.ID))

	if err != nil {
		return false, err
	} else if service != nil {
		return false, nil
	}

	if service.Name != "" {
		service, err = d.findService(ServiceReference(service.Name))

		if err != nil {
			return false, err
		} else if service != nil {
			return false, nil
		}
	}

	return true, nil
}
