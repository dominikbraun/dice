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

// Package storage provides types for persisting entities.
package storage

import "github.com/dominikbraun/dice/entity"

// KVConfig concludes all properties that can be configured by the user.
type KVConfig struct {
	Filename string `json:"filename"`
}

// KV represents a key value store that is aware of all entity stores.
type KV interface {
	NodeStore
	ServiceStore
	InstanceStore
}

// NodeStore provides methods for performing CRUD operations on nodes.
type NodeStore interface {
	CreateNode(node *entity.Node) error
	AllNodes() ([]*entity.Node, error)
	FindNodeBy(identifier interface{}, prop entity.NodeProperty) (*entity.Node, error)
	UpdateNode(identifier interface{}, prop entity.NodeProperty, src entity.Node) error
	DeleteNode(identifier interface{}, prop entity.NodeProperty) error
}

// ServiceStore provides methods for performing CRUD operations on services.
type ServiceStore interface {
	CreateService(node *entity.Service) error
	AllServices() ([]*entity.Service, error)
	FindServiceBy(identifier interface{}, prop entity.ServiceProperty) (*entity.Service, error)
	UpdateService(identifier interface{}, prop entity.ServiceProperty, src entity.Service) error
	DeleteService(identifier interface{}, prop entity.ServiceProperty) error
}

// InstanceStore provides methods for performing CRUD operations on instances.
type InstanceStore interface {
	CreateInstance(node *entity.Instance) error
	AllInstances() ([]*entity.Instance, error)
	FindInstanceBy(identifier interface{}, prop entity.InstanceProperty) (*entity.Instance, error)
	UpdateInstance(identifier interface{}, prop entity.InstanceProperty, src entity.Instance) error
	DeleteInstance(identifier interface{}, prop entity.InstanceProperty) error
}
