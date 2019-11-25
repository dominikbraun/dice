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

type EntityStore interface {
	NodeStore
	ServiceStore
	InstanceStore
}

type NodeStore interface {
	CreateNode(node *entity.Node) error
	FindNodes() ([]*entity.Node, error)
	FindNode(filter NodeFilter) (*entity.Node, error)
	UpdateNode(filter NodeFilter) error
	DeleteNode(filter NodeFilter) error
}

type ServiceStore interface {
	CreateService(service *entity.Service) error
	FindServices() ([]*entity.Service, error)
	FindService(filter ServiceFilter) (entity.Service, error)
	UpdateService(filter ServiceFilter) error
	DeleteService(filter ServiceFilter) error
}

type InstanceStore interface {
	CreateInstance(instance *entity.Instance) error
	FindInstances() ([]*entity.Service, error)
	FindInstance(filter InstanceFilter) (*entity.Service, error)
	UpdateInstance(filter InstanceFilter) error
	DeleteInstance(filter InstanceFilter) error
}
