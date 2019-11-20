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

// Package scheduler provides multiple request scheduling implementations.
package scheduler

import "github.com/dominikbraun/dice/entity"

// BalancingAlgorithm represents a load balancing method.
type BalancingAlgorithm string

const (
	LeastConnectionBalancing    BalancingAlgorithm = "least_connection"
	RandomBalancing             BalancingAlgorithm = "random"
	RoundRobinBalancing         BalancingAlgorithm = "round_robin"
	WeightedRoundRobinBalancing BalancingAlgorithm = "weighted_round_robin"
)

// Scheduler prescribes methods for retrieving the next service instance.
type Scheduler interface {
	Next() (*entity.Instance, error)
}

// Deployment represents a service instance deployed to a node. A scheduler
// needs to know on which node an instance has been deployed, because the
// node's weight is important for some algorithms.
//
// Storing just a reference to the node instance is useful when the node data
// is being updated via the CLI: The new data will be read automatically, since
// the reference points to the global, modified object.
type Deployment struct {
	instance *entity.Instance
	node     *entity.Node
}

// NewScheduler creates a scheduler instance depending on the balancing algorithm.
func NewScheduler(algorithm BalancingAlgorithm, deployments []Deployment) Scheduler {
	var s Scheduler

	switch algorithm {
	case LeastConnectionBalancing:
		panic("unimplemented balancing algorithm")
	case RandomBalancing:
		panic("unimplemented balancing algorithm")
	case RoundRobinBalancing:
		panic("unimplemented balancing algorithm")
	case WeightedRoundRobinBalancing:
		s = newWeightedRoundRobin(deployments)
	default:
		panic("unimplemented balancing algorithm")
	}

	return s
}
