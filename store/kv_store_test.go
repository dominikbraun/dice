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
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/types"
	"net/url"
	"testing"
)

var (
	kvStore *KVStore = nil
)

func setupOnNil(t *testing.T) {
	if kvStore != nil {
		return
	}

	var err error

	kvStore, err = NewKVStore("dice-test-store")
	if err != nil {
		t.Error(err)
	}
}

func TestKVStore_CreateNode(t *testing.T) {
	setupOnNil(t)

	nodeURL, _ := url.Parse("172.21.21.1")
	node, _ := entity.NewNode(nodeURL, types.NodeCreateOptions{Name: "test-node-1"})

	if err := kvStore.CreateNode(node); err != nil {
		t.Error(err.Error())
	}
}

func TestKVStore_FindNode(t *testing.T) {
	setupOnNil(t)

	nodeURL, _ := url.Parse("172.21.21.2")
	node, _ := entity.NewNode(nodeURL, types.NodeCreateOptions{Name: "test-node-2"})

	if err := kvStore.CreateNode(node); err != nil {
		t.Error(err.Error())
	}

	storedNode, err := kvStore.FindNode(node.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if storedNode == nil {
		t.Errorf("node %s is nil", node.ID)
		return
	}

	if storedNode.ID != node.ID {
		t.Errorf("node IDs don't match (%s vs. %s", node.ID, storedNode.ID)
	}
}

func TestKVStore_FindNodes(t *testing.T) {
	setupOnNil(t)

	nodeURL1, _ := url.Parse("172.21.21.3")
	node1, _ := entity.NewNode(nodeURL1, types.NodeCreateOptions{Weight: 255})

	nodeURL2, _ := url.Parse("172.21.21.4")
	node2, _ := entity.NewNode(nodeURL2, types.NodeCreateOptions{Weight: 255})

	if err := kvStore.CreateNode(node1); err != nil {
		t.Error(err)
	}

	if err := kvStore.CreateNode(node2); err != nil {
		t.Error(err)
	}

	nodesByURL, err := kvStore.FindNodes(func(node *entity.Node) bool {
		return node.URL.String() == nodeURL1.String()
	})
	if err != nil {
		t.Error(err)
	}

	if len(nodesByURL) < 1 {
		t.Errorf("%v nodes found, %v expected", len(nodesByURL), 1)
	}

	nodesByWeight, err := kvStore.FindNodes(func(node *entity.Node) bool {
		return node.Weight == 255
	})
	if err != nil {
		t.Error(err)
	}

	if len(nodesByWeight) < 2 {
		t.Errorf("%v nodes found, %v expected", len(nodesByWeight), 2)
	}
}

func TestKVStore_UpdateNode(t *testing.T) {
	setupOnNil(t)

	nodeURL, _ := url.Parse("172.21.21.5")
	node, _ := entity.NewNode(nodeURL, types.NodeCreateOptions{})

	if err := kvStore.CreateNode(node); err != nil {
		t.Error(err)
	}

	node.Weight = 255

	if err := kvStore.UpdateNode(node.ID, node); err != nil {
		t.Error(err)
	}

	updatedNode, err := kvStore.FindNode(node.ID)
	if err != nil {
		t.Error(err)
	}

	if updatedNode.Weight != node.Weight {
		t.Errorf("got weight %v, expected %v", updatedNode.Weight, node.Weight)
	}
}

func TestKVStore_DeleteNode(t *testing.T) {
	setupOnNil(t)

	nodeURL, _ := url.Parse("172.21.21.6")
	node, _ := entity.NewNode(nodeURL, types.NodeCreateOptions{})

	if err := kvStore.CreateNode(node); err != nil {
		t.Error(err)
	}

	if err := kvStore.DeleteNode(node.ID); err != nil {
		t.Error(err)
	}

	deletedNode, err := kvStore.FindNode(node.ID)
	if err != nil {
		t.Error(err)
	}

	if deletedNode != nil {
		t.Errorf("got node %v, expected nil", deletedNode.ID)
	}
}
