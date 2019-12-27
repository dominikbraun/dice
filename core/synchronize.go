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
)

// SynchronizationTask is a type of synchronization between the key-value
// store and the service registry.
type SynchronizationTask uint

const (
	Attach SynchronizationTask = iota
	Detach
	Enable
	Disable
	SetURLs
)

var (
	ErrInvalidSynchronizationTask = errors.New("the provided synchronization task is invalid")
)

// synchronizeNode synchronizes the state of a given node with the state
// of a node that is currently managed by the service registry.
func (d *Dice) synchronizeNode(node *entity.Node, task SynchronizationTask) error {
	switch task {
	case Attach:
		return d.registry.UpdateNodes(func(n *entity.Node) error {
			if n.ID == node.ID {
				n.IsAttached = true
			}
			return nil
		})

	case Detach:
		return d.registry.UpdateNodes(func(n *entity.Node) error {
			if n.ID == node.ID {
				n.IsAttached = false
			}
			return nil
		})

	default:
		return ErrInvalidSynchronizationTask
	}
}

// synchronizeService synchronizes the state of a given service with the
// state of a service that is currently managed by the service registry.
func (d *Dice) synchronizeService(service *entity.Service, task SynchronizationTask) error {
	switch task {
	case Enable:
		return d.registry.UpdateServices(func(s *entity.Service) error {
			if s.ID == service.ID {
				s.IsEnabled = true
			}
			return nil
		})

	case Disable:
		return d.registry.UpdateServices(func(s *entity.Service) error {
			if s.ID == service.ID {
				s.IsEnabled = false
			}
			return nil
		})

	case SetURLs:
		return d.registry.UpdateServices(func(s *entity.Service) error {
			if s.ID == service.ID {
				s.URLs = service.URLs
			}
			return nil
		})

	default:
		return ErrInvalidSynchronizationTask
	}
}

// synchronizeInstance synchronizes the state of a given instance with the
// state of an instance that is currently managed by the service registry.
func (d *Dice) synchronizeInstance(instance *entity.Instance, task SynchronizationTask) error {
	switch task {
	case Attach:
		return d.registry.UpdateInstances(func(i *entity.Instance) error {
			if i.ID == instance.ID {
				i.IsAttached = true
			}
			return nil
		})

	case Detach:
		return d.registry.UpdateInstances(func(i *entity.Instance) error {
			if i.ID == instance.ID {
				i.IsAttached = false
			}
			return nil
		})

	default:
		return ErrInvalidSynchronizationTask
	}
}
