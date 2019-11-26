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

package core

import (
	"errors"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/types"
	"net/url"
)

type NodeReference interface{}

var (
	ErrNodeNotFound = errors.New("node could not be found")
	//ErrNodeAlreadyExists = errors.New("a node with the given URL already exists")
	ErrInvalidNodeReference = errors.New("supplied node identifier is invalid")
)

func (d *Dice) NodeCreate(url *url.URL, options types.NodeCreateOptions) error {
	node, err := entity.NewNode(url, options)
	if err != nil {
		return err
	}

	return d.kvStore.CreateNode(node)
}

func (d *Dice) NodeAttach(nodeRef NodeReference) error {
	node, err := d.getNode(nodeRef)
	if err != nil {
		return err
	}

	if node == nil {
		return ErrNodeNotFound
	}

	node.IsAttached = true

	return d.kvStore.UpdateNode(node.ID, node)
}

func (d *Dice) NodeDetach(nodeRef NodeReference) error {
	node, err := d.getNode(nodeRef)
	if err != nil {
		return err
	}

	if node == nil {
		return ErrNodeNotFound
	}

	node.IsAttached = false

	return d.kvStore.UpdateNode(node.ID, node)
}

func (d *Dice) NodeInfo(nodeRef NodeReference) (types.NodeInfoOutput, error) {
	var nodeInfo types.NodeInfoOutput

	node, err := d.getNode(nodeRef)
	if err != nil {
		return nodeInfo, err
	}

	nodeInfo.Populate(node)

	return nodeInfo, nil
}

func (d *Dice) getNode(nodeRef NodeReference) (*entity.Node, error) {
	ref, ok := nodeRef.(string)
	if !ok {
		return nil, ErrInvalidNodeReference
	}

	nodesByID, err := d.kvStore.FindNodes(func(node *entity.Node) bool {
		return node.ID == ref
	})

	if err != nil {
		return nil, err
	} else if len(nodesByID) > 0 {
		return nodesByID[0], nil
	}

	nodesByName, err := d.kvStore.FindNodes(func(node *entity.Node) bool {
		return node.Name == ref
	})

	if err != nil {
		return nil, err
	} else if len(nodesByName) > 0 {
		return nodesByName[0], nil
	}

	return nil, nil
}
