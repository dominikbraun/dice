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

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/dominikbraun/dice/entity"
)

type Bucket []byte

var (
	diceBucket           Bucket = []byte("dice")
	nodeBucket           Bucket = []byte("nodes")
	serviceBucket        Bucket = []byte("services")
	instanceBucket       Bucket = []byte("instances")
	ErrBucketNotFound    error  = errors.New("bucket could not be found")
	ErrMarshallingFailed error  = errors.New("marshalling of entity failed")
)

type KVStore struct {
	internal *bolt.DB
}

func NewKV(path string) (*KVStore, error) {
	var kv KVStore
	var err error

	if kv.internal, err = bolt.Open(path, 0600, nil); err != nil {
		return nil, err
	}

	if err = (&kv).setup(); err != nil {
		return nil, err
	}

	return &kv, nil
}

func (kv *KVStore) setup() error {
	fn := func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists(diceBucket)
		if err != nil {
			return err
		}

		if _, err = root.CreateBucketIfNotExists(nodeBucket); err != nil {
			return err
		}

		if _, err = root.CreateBucketIfNotExists(serviceBucket); err != nil {
			return err
		}

		if _, err := root.CreateBucketIfNotExists(instanceBucket); err != nil {
			return err
		}

		return nil
	}

	return kv.internal.Update(fn)
}

func (kv *KVStore) set(bucket Bucket, key string, value []byte) error {
	fn := func(tx *bolt.Tx) error {
		b := tx.Bucket(diceBucket).Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}

		return b.Put([]byte(key), value)
	}

	return kv.internal.Update(fn)
}

func (kv *KVStore) getAll(bucket Bucket) ([][]byte, error) {
	var result [][]byte

	fn := func(tx *bolt.Tx) error {
		b := tx.Bucket(diceBucket).Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}

		_ = b.ForEach(func(k, v []byte) error {
			result = append(result, v)
			return nil
		})

		return nil
	}

	if err := kv.internal.View(fn); err != nil {
		return nil, err
	}

	return result, nil
}

func (kv *KVStore) get(bucket Bucket, key string) ([]byte, error) {
	var result []byte

	fn := func(tx *bolt.Tx) error {
		b := tx.Bucket(diceBucket).Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}

		if value := b.Get([]byte(key)); value != nil {
			result = value
			return nil
		}

		return nil
	}

	if err := kv.internal.View(fn); err != nil {
		return nil, err
	}

	return result, nil
}

func (kv *KVStore) delete(bucket Bucket, key string) error {
	fn := func(tx *bolt.Tx) error {
		b := tx.Bucket(diceBucket).Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}

		return b.Delete([]byte(key))
	}

	return kv.internal.Update(fn)
}

func (kv *KVStore) CreateNode(node *entity.Node) error {
	value, err := json.Marshal(node)
	if err != nil {
		return ErrMarshallingFailed
	}

	return kv.set(nodeBucket, node.ID, value)
}

func (kv *KVStore) FindNodes(filter NodeFilter) ([]*entity.Node, error) {
	values, err := kv.getAll(nodeBucket)
	if len(values) == 0 || err != nil {
		return nil, err
	}

	nodes := make([]*entity.Node, len(values))

	for i, v := range values {
		var node entity.Node

		if err = json.Unmarshal(v, &node); err != nil {
			return nil, ErrMarshallingFailed
		}

		if filter != nil && filter(&node) || filter == nil {
			nodes[i] = &node
		}
	}

	return nodes, nil
}

func (kv *KVStore) FindNode(id string) (*entity.Node, error) {
	value, err := kv.get(nodeBucket, id)
	if value == nil || err != nil {
		return nil, err
	}

	var node entity.Node

	if err = json.Unmarshal(value, &node); err != nil {
		return nil, ErrMarshallingFailed
	}

	return &node, nil
}

func (kv *KVStore) UpdateNode(id string, source *entity.Node) error {
	return kv.CreateNode(source)
}

func (kv *KVStore) DeleteNode(id string) error {
	return kv.delete(nodeBucket, id)
}

func (kv *KVStore) CreateService(service *entity.Service) error {
	value, err := json.Marshal(service)
	if err != nil {
		return ErrMarshallingFailed
	}

	return kv.set(serviceBucket, service.ID, value)
}

func (kv *KVStore) FindServices(filter ServiceFilter) ([]*entity.Service, error) {
	values, err := kv.getAll(serviceBucket)
	if len(values) == 0 || err != nil {
		return nil, err
	}

	services := make([]*entity.Service, len(values))

	for i, v := range values {
		var service entity.Service

		if err = json.Unmarshal(v, &service); err != nil {
			return nil, ErrMarshallingFailed
		}

		if filter != nil && filter(&service) || filter == nil {
			services[i] = &service
		}
	}

	return services, nil
}

func (kv *KVStore) FindService(id string) (*entity.Service, error) {
	value, err := kv.get(serviceBucket, id)
	if value == nil || err != nil {
		return nil, err
	}

	var service entity.Service

	if err = json.Unmarshal(value, &service); err != nil {
		return nil, ErrMarshallingFailed
	}

	return &service, nil
}

func (kv *KVStore) UpdateService(id string, source *entity.Service) error {
	return kv.CreateService(source)
}

func (kv *KVStore) DeleteService(id string) error {
	return kv.delete(serviceBucket, id)
}

func (kv *KVStore) CreateInstance(instance *entity.Instance) error {
	value, err := json.Marshal(instance)
	if err != nil {
		return ErrMarshallingFailed
	}

	return kv.set(instanceBucket, instance.ID, value)
}

func (kv *KVStore) FindInstances(filter InstanceFilter) ([]*entity.Instance, error) {
	values, err := kv.getAll(instanceBucket)
	if len(values) == 0 || err != nil {
		return nil, err
	}

	instances := make([]*entity.Instance, len(values))

	for i, v := range values {
		var instance entity.Instance

		if err = json.Unmarshal(v, &instance); err != nil {
			return nil, ErrMarshallingFailed
		}

		if filter != nil && filter(&instance) || filter == nil {
			instances[i] = &instance
		}
	}

	return instances, nil
}

func (kv *KVStore) FindInstance(id string) (*entity.Instance, error) {
	value, err := kv.get(instanceBucket, id)
	if value == nil || err != nil {
		return nil, err
	}

	var instance entity.Instance

	if err = json.Unmarshal(value, &instance); err != nil {
		return nil, ErrMarshallingFailed
	}

	return &instance, nil
}

func (kv *KVStore) UpdateInstance(id string, source *entity.Instance) error {
	return kv.CreateInstance(source)
}

func (kv *KVStore) DeleteInstance(id string) error {
	return kv.delete(instanceBucket, id)
}
