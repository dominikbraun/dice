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
	"github.com/dominikbraun/dice/types"
	"net/url"
)

var (
	ErrInstanceNotFound      = errors.New("instance could not be found")
	ErrInstanceAlreadyExists = errors.New("a instance with the given ID, name or URL already exists")
)

// CreateInstance creates a new instance with the provided service ID,
// node ID and URL. If the `Attach` option is set, the created instance
// will be attached immediately.
func (d *Dice) CreateInstance(serviceID, nodeID string, url *url.URL, options types.InstanceCreateOptions) error {
	instance, err := entity.NewInstance(serviceID, nodeID, url, options)
	if err != nil {
		return err
	}

	isUnique, err := d.instanceIsUnique(instance)

	if err != nil {
		return err
	} else if !isUnique {
		return ErrInstanceAlreadyExists
	}

	if err := d.kvStore.CreateInstance(instance); err != nil {
		return err
	}

	if options.Attach {
		return d.AttachInstance(entity.InstanceReference(instance.ID))
	}

	return nil
}

// AttachInstance attaches an existing instance to Dice, making it available
// as a target for load balancing. This function will update the instance
// data and synchronize the instance with the service registry.
func (d *Dice) AttachInstance(instanceRef entity.InstanceReference) error {
	instance, err := d.findInstance(instanceRef)
	if err != nil {
		return err
	}

	if instance == nil {
		return ErrInstanceNotFound
	}

	instance.IsAttached = true

	if err := d.kvStore.UpdateInstance(instance.ID, instance); err != nil {
		return err
	}

	return d.synchronizeInstance(instance, Attachment)
}

// DetachInstance detaches an existing instance from Dice, removing it as
// a target for load balancing. Detaching an instance will leave all other
// instances of the service untouched.
func (d *Dice) DetachInstance(instanceRef entity.InstanceReference) error {
	instance, err := d.findInstance(instanceRef)
	if err != nil {
		return err
	}

	if instance == nil {
		return ErrInstanceNotFound
	}

	instance.IsAttached = false

	if err := d.kvStore.UpdateInstance(instance.ID, instance); err != nil {
		return err
	}

	return d.synchronizeInstance(instance, Detachment)
}

// InstanceInfo returns user-relevant information for an existing instance.
func (d *Dice) InstanceInfo(instanceRef entity.InstanceReference) (types.InstanceInfoOutput, error) {
	instance, err := d.findInstance(instanceRef)
	if err != nil {
		return types.InstanceInfoOutput{}, err
	}

	instanceInfo := types.InstanceInfoOutput{
		ID:         instance.ID,
		Name:       instance.Name,
		ServiceID:  instance.ServiceID,
		NodeID:     instance.NodeID,
		URL:        instance.URL,
		Version:    instance.Version,
		IsAttached: instance.IsAttached,
		IsAlive:    instance.IsAlive,
	}

	return instanceInfo, nil
}

// findInstance attempts to find an instance in the key-value store that
// matches the reference. The ID has the highest priority, then name and
// URL are checked.
//
// If multiple instances match, only the first one will be returned. If no
// instances match, `nil` - and no error - will be returned.
func (d *Dice) findInstance(instanceRef entity.InstanceReference) (*entity.Instance, error) {
	instancesByID, err := d.kvStore.FindInstances(func(instance *entity.Instance) bool {
		return instance.ID == string(instanceRef)
	})

	if err != nil {
		return nil, err
	} else if len(instancesByID) > 0 {
		return instancesByID[0], nil
	}

	instancesByName, err := d.kvStore.FindInstances(func(instance *entity.Instance) bool {
		return instance.Name == string(instanceRef)
	})

	if err != nil {
		return nil, err
	} else if len(instancesByName) > 0 {
		return instancesByName[0], nil
	}

	instancesByURL, err := d.kvStore.FindInstances(func(instance *entity.Instance) bool {
		return instance.URL.String() == string(instanceRef)
	})

	if err != nil {
		return nil, err
	} else if len(instancesByURL) > 0 {
		return instancesByURL[0], nil
	}

	return nil, nil
}

// instanceIsUnique checks if a newly created instance is unique. An instance
// is unique if no instanceIsUnique with equal identifiers has been found in
// the key value store.
func (d *Dice) instanceIsUnique(instance *entity.Instance) (bool, error) {
	instance, err := d.findInstance(entity.InstanceReference(instance.ID))

	if err != nil {
		return false, err
	} else if instance != nil {
		return false, nil
	}

	instance, err = d.findInstance(entity.InstanceReference(instance.Name))

	if err != nil {
		return false, err
	} else if instance != nil {
		return false, nil
	}

	instance, err = d.findInstance(entity.InstanceReference(instance.URL.String()))

	if err != nil {
		return false, err
	} else if instance != nil {
		return false, nil
	}

	return true, nil
}
