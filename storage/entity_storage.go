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
// Unlike the storage registry, the storage package provides data in a
// representation designed for persisting the data.
package storage

import "github.com/dominikbraun/dice/entity"

// EntityStorage is the common interface for persisting entities.
type EntityStorage interface {
	Create(source entity.Entity, t entity.Type) error
	FindAll(t entity.Type) ([]entity.Entity, error)
	FindBy(identifier interface{}, property entity.Property, t entity.Type) (entity.Entity, error)
	Delete(identifier interface{}, property entity.Property, t entity.Type) error
	Close() error
}
