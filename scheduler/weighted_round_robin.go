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

// Package scheduler provides scheduler implementations for load balancing.
package scheduler

import (
	"errors"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/registry"
)

// WeightedRoundRobin is a scheduler that basically implements the Round
// Robin algorithm under consideration of node weights.
//
// Taking the node weights into account means that a node of weight 2 will
// be selected twice as often as a node of weight 1 - more exactly, the
// instance deployed to that node will be selected. If there are two service
// instances deployed to a node of weight 2, the node will receive four times
// more requests as a consequence.
//
// Instances that are either detached or considered dead won't be selected,
// just as instances that are deployed to a detached or dead node.
type WeightedRoundRobin struct {
	deployments   *[]registry.Deployment
	currentIndex  int
	currentWeight uint8
}

// newWeightedRoundRobin creates a new WeightedRoundRobin instance.
func newWeightedRoundRobin(deployments *[]registry.Deployment) *WeightedRoundRobin {
	wrr := WeightedRoundRobin{
		deployments:   deployments,
		currentIndex:  0,
		currentWeight: uint8(0),
	}

	return &wrr
}

// Next implements registry.Scheduler.Next. It is an implementation of the
// Weighted Round Robin algorithm, respecting the rules described above.
func (wrr *WeightedRoundRobin) Next() (*entity.Instance, error) {
	attempts := 0

lookup:
	for attempts < len(*wrr.deployments) {
		index := wrr.currentIndex % len(*wrr.deployments)
		d := (*wrr.deployments)[index]

		if !d.Instance.IsAttached || !d.Instance.IsAlive {
			wrr.currentIndex++
			attempts++
			continue lookup
		}

		if d.Node.Weight == wrr.currentWeight {
			wrr.currentIndex++
			wrr.currentWeight = uint8(0)
			attempts++
			continue lookup
		}

		if d.Node.Weight > wrr.currentWeight {
			wrr.currentWeight++
			return d.Instance, nil
		}

		attempts++
	}

	return nil, errors.New("no service instance found")
}
