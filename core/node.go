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
	"github.com/dominikbraun/dice/types"
	"net/url"
)

// NodeReference is a string that identifies a node, e. g. an ID or name.
type NodeReference string

var (
	ErrNodeNotFound      = errors.New("node could not be found")
	ErrNodeAlreadyExists = errors.New("a node with the given ID already exists")
)

// CreateNode creates a new node with the provided URL and stores the node
// in the key-value store. If the `Attach` option is set, the created node
// will be attached immediately.
func (d *Dice) CreateNode(url *url.URL, options types.NodeCreateOptions) error {
	node, err := entity.NewNode(url, options)
	if err != nil {
		return err
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
		return d.AttachNode(NodeReference(node.ID))
	}

	return nil
}

// AttachNode attaches an existing node to Dice, making it available as a
// target for load balancing. This function will update the node data and
// synchronize the node with the service registry.
func (d *Dice) AttachNode(nodeRef NodeReference) error {
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

	return d.synchronizeNode(node, Attachment)
}

// DetachNode detaches an existing node from Dice, removing it as a target
// for load balancing. Detaching a node will make all instances deployed to
// that node unavailable until it gets attached again.
func (d *Dice) DetachNode(nodeRef NodeReference) error {
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

	return d.synchronizeNode(node, Detachment)
}

// NodeInfo returns user-relevant information for an existing node.
func (d *Dice) NodeInfo(nodeRef NodeReference) (types.NodeInfoOutput, error) {
	node, err := d.findNode(nodeRef)
	if err != nil {
		return types.NodeInfoOutput{}, err
	}

	nodeInfo := types.NodeInfoOutput{
		ID:         node.ID,
		Name:       node.Name,
		URL:        node.URL,
		IsAttached: node.IsAttached,
		IsAlive:    node.IsAlive,
	}

	return nodeInfo, nil
}

// findNode attempts to find a node in the key-value store that matches the
// reference. The ID has the highest priority, then name and URL are checked.
//
// If multiple nodes match, only the first one will be returned. If no nodes
// match, `nil` - and no error - will be returned.
func (d *Dice) findNode(nodeRef NodeReference) (*entity.Node, error) {
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
	node, err := d.findNode(NodeReference(node.ID))

	if err != nil {
		return false, err
	} else if node != nil {
		return false, nil
	}

	node, err = d.findNode(NodeReference(node.Name))

	if err != nil {
		return false, err
	} else if node != nil {
		return false, nil
	}

	node, err = d.findNode(NodeReference(node.URL.String()))

	if err != nil {
		return false, err
	} else if node != nil {
		return false, nil
	}

	return true, nil
}
