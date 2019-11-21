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

// Entity is any Dice core entity that will be stored.
type Entity interface{}

// Property is an entity's property it will be identified by.
type Property interface{}

// EntityType indicates the data type of the Entity. Depending on the type,
// an EntityStore has to decide where and how the entity has to be stored.
type EntityType uint

const (
	Node     EntityType = 0
	Service  EntityType = 1
	Instance EntityType = 2
)

// EntityStore is the common interface for persisting core entities.
type EntityStore interface {
	Create(entity Entity, t EntityType) error
	FindAll(t EntityType) ([]Entity, error)
	FindBy(identifier string, property Property, t EntityType) (Entity, error)
	Update(identifier string, property Property, t EntityType, src Entity) error
	Delete(identifier string, property Property, t EntityType) error
}
