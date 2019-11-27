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
	"github.com/dominikbraun/dice/registry"
)

type SynchronizationTask uint

const (
	Attachment SynchronizationTask = 0
	Detachment SynchronizationTask = 1
	Removal    SynchronizationTask = 2
)

var (
	ErrInvalidSynchronizationTask = errors.New("the provided synchronization task is invalid")
)

func (d *Dice) synchronizeNode(node *entity.Node, task SynchronizationTask) error {
	switch task {
	case Attachment:
		filter := func(n *entity.Node) bool {
			return n.ID == node.ID
		}
		update := func(n *entity.Node) error {
			n.IsAttached = true
			return nil
		}
		return d.registry.UpdateNodes(filter, update)

	case Detachment:
		filter := func(n *entity.Node) bool {
			return n.ID == node.ID
		}
		update := func(n *entity.Node) error {
			n.IsAttached = false
			return nil
		}
		return d.registry.UpdateNodes(filter, update)

	case Removal:
		filter := func(deployment registry.Deployment) bool {
			return deployment.Node.ID == node.ID
		}
		update := func(deployment registry.Deployment) error {
			deployment.IsValid = false
			return nil
		}
		return d.registry.UpdateDeployments(filter, update)

	default:
		return ErrInvalidSynchronizationTask
	}
}
