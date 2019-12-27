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

// Package core provides the Dice load balancer and its methods.
package core

import (
	"errors"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/store"
	"github.com/dominikbraun/dice/types"
	"net/url"
)

var (
	ErrNodeNotFound      = errors.New("node could not be found")
	ErrNodeAlreadyExists = errors.New("the given node already exists")
)

// CreateNode creates a new node with the provided URL and stores the node
// in the key-value store. If the `Attach` option is set, the created node
// will be attached immediately.
func (d *Dice) CreateNode(url *url.URL, options types.NodeCreateOptions) error {
	node, err := entity.NewNode(url, options)
	if err != nil {
		return err
	}

	if ok, message := validateNode(node); !ok {
		return errors.New(message)
	}

	isUnique, err := d.nodeIsUnique(node)

	if err != nil {
		return err
	} else if !isUnique {
		return ErrNodeAlreadyExists
	}

	if err := d.kvStore.CreateNode(node); err != nil {
		return err
	}

	if options.Attach {
		return d.AttachNode(entity.NodeReference(node.ID))
	}

	return nil
}

// AttachNode attaches an existing node to Dice, making it available as a
// target for load balancing. This function will update the node data and
// synchronize the node with the service registry.
func (d *Dice) AttachNode(nodeRef entity.NodeReference) error {
	node, err := d.findNode(nodeRef)
	if err != nil {
		return err
	}

	if node == nil {
		return ErrNodeNotFound
	}

	node.IsAttached = true

	if err := d.kvStore.UpdateNode(node.ID, node); err != nil {
		return err
	}

	return d.synchronizeNode(node, Attach)
}

// DetachNode detaches an existing node from Dice, removing it as a target
// for load balancing. Detaching a node will make all instances deployed to
// that node unavailable until it gets attached again.
func (d *Dice) DetachNode(nodeRef entity.NodeReference) error {
	node, err := d.findNode(nodeRef)
	if err != nil {
		return err
	}

	if node == nil {
		return ErrNodeNotFound
	}

	node.IsAttached = false

	if err := d.kvStore.UpdateNode(node.ID, node); err != nil {
		return err
	}

	return d.synchronizeNode(node, Detach)
}

// NodeInfo returns user-relevant information for an existing node.
func (d *Dice) NodeInfo(nodeRef entity.NodeReference) (types.NodeInfoOutput, error) {
	node, err := d.findNode(nodeRef)
	if err != nil {
		return types.NodeInfoOutput{}, err
	}

	if node == nil {
		return types.NodeInfoOutput{}, ErrNodeNotFound
	}

	nodeInfo := types.NodeInfoOutput{
		ID:         node.ID,
		Name:       node.Name,
		URL:        node.URL.String(),
		IsAttached: node.IsAttached,
		IsAlive:    node.IsAlive,
	}

	return nodeInfo, nil
}

// ListNodes returns a list of stored nodes. By default, detached nodes will
// be ignored. They only will be returned if the options say to do so. In any
// case, dead nodes will be returned.
func (d *Dice) ListNodes(options types.NodeListOptions) ([]types.NodeInfoOutput, error) {
	filter := store.AllNodesFilter

	if !options.All {
		filter = func(node *entity.Node) bool {
			return node.IsAttached
		}
	}

	nodes, err := d.kvStore.FindNodes(filter)
	if err != nil {
		return nil, err
	}

	nodeList := make([]types.NodeInfoOutput, len(nodes))

	for i, n := range nodes {
		info := types.NodeInfoOutput{
			ID:         n.ID,
			Name:       n.Name,
			URL:        n.URL.String(),
			IsAttached: n.IsAttached,
			IsAlive:    n.IsAlive,
		}
		nodeList[i] = info
	}

	return nodeList, nil
}

// findNode attempts to find a node in the key-value store that matches the
// reference. The ID has the highest priority, then name and URL are checked.
//
// If multiple nodes match, only the first one will be returned. If no nodes
// match, `nil` - and no error - will be returned.
func (d *Dice) findNode(nodeRef entity.NodeReference) (*entity.Node, error) {
	nodesByID, err := d.kvStore.FindNodes(func(node *entity.Node) bool {
		return node.ID == string(nodeRef)
	})

	if err != nil {
		return nil, err
	} else if len(nodesByID) > 0 {
		return nodesByID[0], nil
	}

	nodesByName, err := d.kvStore.FindNodes(func(node *entity.Node) bool {
		return node.Name == string(nodeRef)
	})

	if err != nil {
		return nil, err
	} else if len(nodesByName) > 0 {
		return nodesByName[0], nil
	}

	nodesByURL, err := d.kvStore.FindNodes(func(node *entity.Node) bool {
		return node.URL.String() == string(nodeRef)
	})

	if err != nil {
		return nil, err
	} else if len(nodesByURL) > 0 {
		return nodesByURL[0], nil
	}

	return nil, nil
}

// nodeIsUnique checks if a newly created node is unique. A node is unique
// if no node with equal identifiers has been found in the key value store.
func (d *Dice) nodeIsUnique(node *entity.Node) (bool, error) {
	storedNode, err := d.findNode(entity.NodeReference(node.ID))

	if err != nil {
		return false, err
	} else if storedNode != nil {
		return false, nil
	}

	if node.Name != "" {
		storedNode, err = d.findNode(entity.NodeReference(node.Name))

		if err != nil {
			return false, err
		} else if storedNode != nil {
			return false, nil
		}
	}

	storedNode, err = d.findNode(entity.NodeReference(node.URL.String()))

	if err != nil {
		return false, err
	} else if storedNode != nil {
		return false, nil
	}

	return true, nil
}
