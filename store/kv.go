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

type KV struct {
	internal *bolt.DB
}

func NewKV(path string) (*KV, error) {
	var kv KV
	var err error

	if kv.internal, err = bolt.Open(path, 0600, nil); err != nil {
		return nil, err
	}

	if err = (&kv).setup(); err != nil {
		return nil, err
	}

	return &kv, nil
}

func (kv *KV) setup() error {
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

func (kv *KV) set(bucket Bucket, key string, value []byte) error {
	fn := func(tx *bolt.Tx) error {
		b := tx.Bucket(diceBucket).Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}

		return b.Put([]byte(key), value)
	}

	return kv.internal.Update(fn)
}

func (kv *KV) getAll(bucket Bucket) ([][]byte, error) {
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

func (kv *KV) get(bucket Bucket, key string) ([]byte, error) {
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

func (kv *KV) Delete(bucket Bucket, key string) error {
	fn := func(tx *bolt.Tx) error {
		b := tx.Bucket(diceBucket).Bucket(bucket)
		if b == nil {
			return ErrBucketNotFound
		}

		return b.Delete([]byte(key))
	}

	return kv.internal.Update(fn)
}

func (kv *KV) CreateNode(node *entity.Node) error {
	value, err := json.Marshal(node)
	if err != nil {
		return ErrMarshallingFailed
	}

	return kv.set(nodeBucket, node.ID, value)
}

func (kv *KV) FindNodes(filter NodeFilter) ([]*entity.Node, error) {
	values, err := kv.getAll(nodeBucket)
	if err != nil {
		return nil, err
	}

	nodes := make([]*entity.Node, len(values))

	for _, v := range values {
		var node *entity.Node

		if err = json.Unmarshal(v, &node); err != nil {
			return nil, ErrMarshallingFailed
		}

		if filter != nil && filter(node) || filter == nil {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func (kv *KV) FindNode(id string) (*entity.Node, error) {
	value, err := kv.get(nodeBucket, id)
	if err != nil {
		return nil, err
	}

	var node *entity.Node

	if err = json.Unmarshal(value, node); err != nil {
		return nil, ErrMarshallingFailed
	}

	return node, nil
}

func (kv *KV) UpdateNode(id string, source *entity.Node) error {
	return kv.CreateNode(source)
}

func (kv *KV) DeleteNode(id string) error {
	return kv.Delete(nodeBucket, id)
}

func (kv *KV) CreateService(service *entity.Service) error {
	value, err := json.Marshal(service)
	if err != nil {
		return ErrMarshallingFailed
	}

	return kv.set(serviceBucket, service.ID, value)
}

func (kv *KV) FindServices(filter ServiceFilter) ([]*entity.Service, error) {
	values, err := kv.getAll(serviceBucket)
	if err != nil {
		return nil, err
	}

	services := make([]*entity.Service, len(values))

	for _, v := range values {
		var service *entity.Service

		if err = json.Unmarshal(v, &service); err != nil {
			return nil, ErrMarshallingFailed
		}

		if filter != nil && filter(service) || filter == nil {
			services = append(services, service)
		}
	}

	return services, nil
}

func (kv *KV) FindService(id string) (*entity.Service, error) {
	value, err := kv.get(serviceBucket, id)
	if err != nil {
		return nil, err
	}

	var service *entity.Service

	if err = json.Unmarshal(value, service); err != nil {
		return nil, ErrMarshallingFailed
	}

	return service, nil
}

func (kv *KV) UpdateService(id string, source *entity.Service) error {
	return kv.CreateService(source)
}

func (kv *KV) DeleteService(id string) error {
	return kv.Delete(serviceBucket, id)
}

func (kv *KV) CreateInstance(instance *entity.Instance) error {
	value, err := json.Marshal(instance)
	if err != nil {
		return ErrMarshallingFailed
	}

	return kv.set(instanceBucket, instance.ID, value)
}

func (kv *KV) FindInstances(filter InstanceFilter) ([]*entity.Instance, error) {
	values, err := kv.getAll(instanceBucket)
	if err != nil {
		return nil, err
	}

	instances := make([]*entity.Instance, len(values))

	for _, v := range values {
		var instance *entity.Instance

		if err = json.Unmarshal(v, &instance); err != nil {
			return nil, ErrMarshallingFailed
		}

		if filter != nil && filter(instance) || filter == nil {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func (kv *KV) FindInstance(id string) (*entity.Instance, error) {
	value, err := kv.get(instanceBucket, id)
	if err != nil {
		return nil, err
	}

	var instance *entity.Instance

	if err = json.Unmarshal(value, instance); err != nil {
		return nil, ErrMarshallingFailed
	}

	return instance, nil
}

func (kv *KV) UpdateInstance(id string, source *entity.Instance) error {
	return kv.CreateInstance(source)
}

func (kv *KV) DeleteInstance(id string) error {
	return kv.Delete(instanceBucket, id)
}
