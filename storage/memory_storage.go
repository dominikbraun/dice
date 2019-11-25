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

// Package storage provides persistence interfaces and implementations.
package storage

import (
	"errors"
	"fmt"

	"github.com/dominikbraun/dice/entity"
)

// MemoryStorage represents a simple in-memory storage. Manipulating a stored
// entity will take effect on any function reading the entity.
type MemoryStorage struct {
	nodes     []*entity.Node
	services  []*entity.Service
	instances []*entity.Instance
}

// NewMemoryStorage creates a new Memory instances that will be initialized with the
// pre-allocated entity slices.
func NewMemoryStorage(nodes []*entity.Node, services []*entity.Service, instances []*entity.Instance) *MemoryStorage {
	m := MemoryStorage{
		nodes:     nodes,
		services:  services,
		instances: instances,
	}

	return &m
}

// Create implements Entity.Create.
func (m *MemoryStorage) Create(source entity.Entity, t entity.Type) error {
	switch t {
	case entity.TypeNode:
		node, ok := source.(*entity.Node)
		if !ok {
			return typeAssertionErr("*entity.Node")
		}
		m.nodes = append(m.nodes, node)

	case entity.TypeService:
		service, ok := source.(*entity.Service)
		if !ok {
			return typeAssertionErr("*entity.Service")
		}
		m.services = append(m.services, service)

	case entity.TypeInstance:
		instance, ok := source.(*entity.Instance)
		if !ok {
			return typeAssertionErr("*entity.Instance")
		}
		m.instances = append(m.instances, instance)

	default:
		return invalidEntityTypeErr()
	}

	return nil
}

// FindAll implements Entity.FindAll.
func (m *MemoryStorage) FindAll(t entity.Type) ([]entity.Entity, error) {
	switch t {
	case entity.TypeNode:
		nodes := make([]entity.Entity, len(m.nodes))

		for i, n := range m.nodes {
			nodes[i] = n
		}

		return nodes, nil

	case entity.TypeService:
		services := make([]entity.Entity, len(m.services))

		for i, s := range m.services {
			services[i] = s
		}

		return services, nil

	case entity.TypeInstance:
		instances := make([]entity.Entity, len(m.instances))

		for i, inst := range m.services {
			instances[i] = inst
		}

		return instances, nil

	default:
		return nil, invalidEntityTypeErr()
	}
}

// FindBy implements Entity.FindBy.
func (m *MemoryStorage) FindBy(identifier interface{}, property entity.Property, t entity.Type) (entity.Entity, error) {
	switch t {
	case entity.TypeNode:
		for _, n := range m.nodes {
			if property == entity.NodeID && identifier == n.ID {
				return n, nil
			}
			if property == entity.NodeName && identifier == n.Config.Name {
				return n, nil
			}
			if property == entity.NodeURL && identifier == n.Config.URL {
				return n, nil
			}
		}

	case entity.TypeService:
		for _, s := range m.services {
			if property == entity.ServiceID && identifier == s.ID {
				return s, nil
			}
			if property == entity.ServiceName && identifier == s.Config.Name {
				return s, nil
			}
		}

	case entity.TypeInstance:
		for _, i := range m.instances {
			if property == entity.InstanceID && identifier == i.ID {
				return i, nil
			}
			if property == entity.InstanceName && identifier == i.Config.Name {
				return i, nil
			}
			if property == entity.InstanceURL && identifier == i.Config.URL {
				return i, nil
			}
		}
	}

	return nil, invalidEntityTypeErr()
}

// Delete implements Entity.Delete.
func (m *MemoryStorage) Delete(identifier interface{}, property entity.Property, t entity.Type) error {
	switch t {
	case entity.TypeNode:
		indexOf := -1

		for i, n := range m.nodes {
			if property == entity.NodeID && identifier == n.ID {
				indexOf = i
			}
			if property == entity.NodeName && identifier == n.Config.Name {
				indexOf = i
			}
			if property == entity.NodeURL && identifier == n.Config.URL {
				indexOf = i
			}
		}

		if indexOf != -1 {
			m.nodes[indexOf] = m.nodes[len(m.nodes)-1]
			m.nodes = m.nodes[:len(m.nodes)-1]
		} else {
			return entityNotFoundErr(identifier)
		}

	case entity.TypeService:
		indexOf := -1

		for i, s := range m.services {
			if property == entity.ServiceID && identifier == s.ID {
				indexOf = i
			}
			if property == entity.ServiceName && identifier == s.Config.Name {
				indexOf = i
			}
		}

		if indexOf != -1 {
			m.services[indexOf] = m.services[len(m.services)-1]
			m.services = m.services[:len(m.services)-1]
		} else {
			return entityNotFoundErr(identifier)
		}

	case entity.TypeInstance:
		indexOf := -1

		for i, inst := range m.instances {
			if property == entity.InstanceID && identifier == inst.ID {
				indexOf = i
			}
			if property == entity.InstanceName && identifier == inst.Config.Name {
				indexOf = i
			}
			if property == entity.InstanceURL && identifier == inst.Config.URL {
				indexOf = i
			}
		}

		if indexOf != -1 {
			m.instances[indexOf] = m.instances[len(m.instances)-1]
			m.instances = m.instances[:len(m.instances)-1]
		} else {
			return entityNotFoundErr(identifier)
		}

	default:
		return invalidEntityTypeErr()
	}

	return nil
}

// Close implements EntityStorage.Close.
func (m *MemoryStorage) Close() error {
	return nil
}

// typeAssertionErr returns an error indicating that a type assertion has failed.
func typeAssertionErr(asserted string) error {
	err := fmt.Errorf("entity is not of type %v", asserted)
	return err
}

// entityNotFoundErr returns an error indicating that a specific entity could not
// be found. Will not be used by Find and FindAll.
func entityNotFoundErr(identifier interface{}) error {
	err := fmt.Errorf("%v could not be found", identifier)
	return err
}

// invalidEntityTypeErr returns an error indicating that a EntityType is invalid.
func invalidEntityTypeErr() error {
	err := errors.New("unsupported entity type")
	return err
}
