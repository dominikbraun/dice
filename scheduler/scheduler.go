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
	"github.com/dominikbraun/dice/registry"
)

// BalancingMethod describes a load balancing algorithm.
type BalancingMethod string

const (
	LeastConnectionBalancing    BalancingMethod = "least_connection"
	RandomBalancing             BalancingMethod = "random"
	RoundRobinBalancing         BalancingMethod = "round_robin"
	WeightedRoundRobinBalancing BalancingMethod = "weighted_round_robin"
)

var (
	ErrUnsupportedMethod = errors.New("balancing method is not supported")
)

// New creates a new Scheduler instance depending on the provided balancing
// method. The particular instance has read-only access to the deployments.
func New(deployments *[]registry.Deployment, method BalancingMethod) (registry.Scheduler, error) {
	switch method {
	case WeightedRoundRobinBalancing:
		return newWeightedRoundRobin(deployments), nil
	default:
		return nil, ErrUnsupportedMethod
	}
}
