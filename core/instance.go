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
	"fmt"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/registry"
	"github.com/dominikbraun/dice/store"
	"github.com/dominikbraun/dice/types"
	"strings"
)

var (
	ErrInstanceNotFound      = errors.New("instance could not be found")
	ErrInstanceAlreadyExists = errors.New("a instance with the given ID, name or URL already exists")
)

// CreateInstance creates a new instance with the provided service ID, node
// ID and port. If the `Attach` option is set, the created instance will be
// attached immediately.
func (d *Dice) CreateInstance(serviceRef entity.ServiceReference, nodeRef entity.NodeReference, url string, options types.InstanceCreateOptions) error {
	service, err := d.findService(serviceRef)

	if err != nil {
		return err
	} else if service == nil {
		return ErrServiceNotFound
	}

	node, err := d.findNode(nodeRef)

	if err != nil {
		return err
	} else if node == nil {
		return ErrNodeNotFound
	}

	instance, err := entity.NewInstance(service.ID, node.ID, normalizeURL(url), options)
	if err != nil {
		return err
	}

	if ok, message := validateInstance(instance); !ok {
		return errors.New(message)
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

	deployment := registry.Deployment{
		Node:     node,
		Instance: instance,
	}

	if err := d.registry.RegisterDeployment(deployment); err != nil {
		return err
	}

	if options.Attach {
		if err := d.AttachInstance(entity.InstanceReference(instance.ID)); err != nil {
			return fmt.Errorf("instance created but not attached: %s", err.Error())
		}
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
	} else if instance == nil {
		return ErrInstanceNotFound
	}

	instance.IsAttached = true

	if err := d.kvStore.UpdateInstance(instance.ID, instance); err != nil {
		return err
	}

	return d.registry.Update(func(s *registry.Service) error {
		for _, d := range s.Deployments {
			if d.Instance.ID == instance.ID {
				d.Instance.IsAttached = true
			}
		}
		return nil
	})
}

// DetachInstance detaches an existing instance from Dice, removing it as
// a target for load balancing. Detaching an instance will leave all other
// instances of the service untouched.
func (d *Dice) DetachInstance(instanceRef entity.InstanceReference) error {
	instance, err := d.findInstance(instanceRef)

	if err != nil {
		return err
	} else if instance == nil {
		return ErrInstanceNotFound
	}

	instance.IsAttached = false

	if err := d.kvStore.UpdateInstance(instance.ID, instance); err != nil {
		return err
	}

	return d.registry.Update(func(s *registry.Service) error {
		for _, d := range s.Deployments {
			if d.Instance.ID == instance.ID {
				d.Instance.IsAttached = false
			}
		}
		return nil
	})
}

// RemoveInstance removes an instance entirely. After getting unregistered
// from the service registry, it won't be available for load balancing any
// longer. Also, it can't be restored anymore.
//
// Returns an error in case the instance can't be removed safely, unless
// --force is used.
func (d *Dice) RemoveInstance(instanceRef entity.InstanceReference, options types.InstanceRemoveOptions) error {
	instance, err := d.findInstance(instanceRef)

	if err != nil {
		return err
	} else if instance == nil {
		return ErrInstanceNotFound
	}

	filter := func(deployment registry.Deployment) bool {
		return deployment.Instance.ID == instance.ID
	}

	if ok := d.registry.UnregisterDeployments(filter, options.Force); !ok {
		return fmt.Errorf("instance is attached, detach it or use --force")
	}

	return d.kvStore.DeleteInstance(instance.ID)
}

// InstanceInfo returns user-relevant information for an existing instance.
func (d *Dice) InstanceInfo(instanceRef entity.InstanceReference) (types.InstanceInfoOutput, error) {
	instance, err := d.findInstance(instanceRef)

	if err != nil {
		return types.InstanceInfoOutput{}, err
	} else if instance == nil {
		return types.InstanceInfoOutput{}, ErrInstanceNotFound
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

// ListInstances returns a list of stored instances. By default, detached
// instances will be ignored. They only will be returned if the options say
// to do so.
func (d *Dice) ListInstances(options types.InstanceListOptions) ([]types.InstanceInfoOutput, error) {
	filter := store.AllInstancesFilter

	if !options.All {
		filter = func(instance *entity.Instance) bool {
			return instance.IsAttached
		}
	}

	instances, err := d.kvStore.FindInstances(filter)
	if err != nil {
		return nil, err
	}

	serviceList := make([]types.InstanceInfoOutput, len(instances))

	for i, inst := range instances {
		info := types.InstanceInfoOutput{
			ID:         inst.ID,
			Name:       inst.Name,
			ServiceID:  inst.ServiceID,
			NodeID:     inst.NodeID,
			URL:        inst.URL,
			Version:    inst.Version,
			IsAttached: inst.IsAttached,
			IsAlive:    inst.IsAlive,
		}
		serviceList[i] = info
	}

	return serviceList, nil
}

// findInstance attempts to find an instance in the key-value store that
// matches the reference. The ID has the highest priority, then name and
// URL are checked.
//
// In order to identify the instance by its URL, a node with the provided
// URL will be searched. If an instance with the URL's port is deployed to
// that node, that instance will be selected.
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

	instanceURL := normalizeURL(string(instanceRef))

	instancesByURL, err := d.kvStore.FindInstances(func(i *entity.Instance) bool {
		return i.URL == instanceURL
	})

	if err != nil {
		return nil, err
	} else if len(instanceURL) > 0 {
		return instancesByURL[0], nil
	}

	return nil, nil
}

// instanceIsUnique checks if a newly created instance is unique. An instance
// is unique if no instanceIsUnique with equal identifiers has been found in
// the key value store.
func (d *Dice) instanceIsUnique(instance *entity.Instance) (bool, error) {
	instancesByURL, err := d.kvStore.FindInstances(func(i *entity.Instance) bool {
		return i.URL == instance.URL
	})

	if err != nil {
		return false, err
	} else if len(instancesByURL) > 0 {
		return false, nil
	}

	if instance.Name != "" {
		instancesByName, err := d.kvStore.FindInstances(func(i *entity.Instance) bool {
			return i.ServiceID == instance.ServiceID && i.Name == instance.Name
		})

		if err != nil {
			return false, err
		} else if len(instancesByName) > 0 {
			return false, nil
		}
	}

	return true, nil
}

// normalizeURL turns any URL string into an normalized, uniformly URL. This
// is necessary for converting a user input like example.com into an appropriate
// url.URL instance.
//
// Even though example.com is a valid URL for url.Parse(), it is not possible to
// dial it since the scheme is missing. Only //example.com would be usable, and
// normalizeURL makes sure that the provided URL will be usable.
func normalizeURL(url string) string {
	normalized := url

	if strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "//") {
		normalized = "//" + normalized
	}

	return normalized
}
