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

package store

import "github.com/dominikbraun/dice/entity"

type (
	NodeFilter     func(node *entity.Node) bool
	ServiceFilter  func(service *entity.Service) bool
	InstanceFilter func(instance *entity.Instance) bool
)

var (
	AllNodesFilter     NodeFilter     = func(node *entity.Node) bool { return true }
	AllServicesFilter  ServiceFilter  = func(service *entity.Service) bool { return true }
	AllInstancesFilter InstanceFilter = func(instance *entity.Instance) bool { return true }
)

type EntityStore interface {
	NodeStore
	ServiceStore
	InstanceStore
	Close() error
}

type NodeStore interface {
	CreateNode(node *entity.Node) error
	FindNodes(filter NodeFilter) ([]*entity.Node, error)
	FindNode(id string) (*entity.Node, error)
	UpdateNode(id string, source *entity.Node) error
	DeleteNode(id string) error
}

type ServiceStore interface {
	CreateService(service *entity.Service) error
	FindServices(filter ServiceFilter) ([]*entity.Service, error)
	FindService(id string) (*entity.Service, error)
	UpdateService(id string, source *entity.Service) error
	DeleteService(id string) error
}

type InstanceStore interface {
	CreateInstance(instance *entity.Instance) error
	FindInstances(filter InstanceFilter) ([]*entity.Instance, error)
	FindInstance(id string) (*entity.Instance, error)
	UpdateInstance(id string, source *entity.Instance) error
	DeleteInstance(id string) error
}
