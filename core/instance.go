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
	"net/url"
)

type InstanceReference string

var (
	ErrInstanceNotFound      = errors.New("instance could not be found")
	ErrInstanceAlreadyExists = errors.New("a instance with the given ID, name or URL already exists")
)

func (d *Dice) InstanceCreate(serviceID, nodeID string, url *url.URL, options types.InstanceCreateOptions) error {
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
		return d.InstanceAttach(InstanceReference(instance.ID))
	}

	return nil
}

func (d *Dice) InstanceAttach(instanceRef InstanceReference) error {
	instance, err := d.findInstance(instanceRef)
	if err != nil {
		return err
	}

	if instance == nil {
		return ErrInstanceNotFound
	}

	instance.IsAttached = true

	return d.kvStore.UpdateInstance(instance.ID, instance)
}

func (d *Dice) InstanceDetach(instanceRef InstanceReference) error {
	instance, err := d.findInstance(instanceRef)
	if err != nil {
		return err
	}

	if instance == nil {
		return ErrInstanceNotFound
	}

	instance.IsAttached = false

	return d.kvStore.UpdateInstance(instance.ID, instance)
}

func (d *Dice) InstanceInfo(instanceRef InstanceReference) (types.InstanceInfoOutput, error) {
	var instanceInfo types.InstanceInfoOutput

	instance, err := d.findInstance(instanceRef)
	if err != nil {
		return instanceInfo, err
	}

	instanceInfo.Populate(instance)

	return instanceInfo, nil
}

func (d *Dice) findInstance(instanceRef InstanceReference) (*entity.Instance, error) {
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

func (d *Dice) instanceIsUnique(instance *entity.Instance) (bool, error) {
	instance, err := d.findInstance(InstanceReference(instance.ID))

	if err != nil {
		return false, err
	} else if instance != nil {
		return false, nil
	}

	instance, err = d.findInstance(InstanceReference(instance.Name))

	if err != nil {
		return false, err
	} else if instance != nil {
		return false, nil
	}

	instance, err = d.findInstance(InstanceReference(instance.URL.String()))

	if err != nil {
		return false, err
	} else if instance != nil {
		return false, nil
	}

	return true, nil
}
