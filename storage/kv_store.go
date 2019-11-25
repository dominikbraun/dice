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

import "github.com/dgraph-io/badger"

import "github.com/dominikbraun/dice/entity"

const (
	rootBucket     = "dice"
	nodeBucket     = "nodes"
	serviceBucket  = "services"
	instanceBucket = "instances"
)

// KVStoreConfig concludes all properties that can be configured by the user.
type KVStoreConfig struct {
	Path string `json:"path"`
}

// KVStore is a simple key value store that persists the created an modified
// entities. Its entries will be loaded into the memory storage on startup.
type KVStore struct {
	config   KVStoreConfig
	internal *badger.DB
}

// NewKVStore creates a new KVStore instance and opens the internal database.
func NewKVStore(config KVStoreConfig) (*KVStore, error) {
	kv := KVStore{
		config: config,
	}

	options := badger.DefaultOptions(kv.config.Path)
	var err error

	if kv.internal, err = badger.Open(options); err != nil {
		return nil, err
	}

	return &kv, nil
}

// Create implements EntityStorage.Create.
func (kv *KVStore) Create(source entity.Entity, t entity.Type) error {
	return nil
}

// FindAll implements EntityStorage.FindAll.
func (kv *KVStore) FindAll(t entity.Type) ([]entity.Entity, error) {
	return nil, nil
}

// FindBy implements EntityStorage.FindBy.
func (kv *KVStore) FindBy(identifier interface{}, property entity.Property, t entity.Type) ([]entity.Entity, error) {
	return nil, nil
}

// Delete implements EntityStorage.Delete.
func (kv *KVStore) Delete(identifier interface{}, property entity.Property, t entity.Type) error {
	return nil
}

// Close implements EntityStorage.Close.
func (kv *KVStore) Close() error {
	return kv.internal.Close()
}
