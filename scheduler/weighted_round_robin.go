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

import (
	"errors"

	"github.com/dominikbraun/dice/entity"
)

// WeightedRoundRobin is a request scheduler which implements the Weighted
// Round Robin algorithm.
type WeightedRoundRobin struct {
	deployments   *[]entity.Deployment
	currentIndex  int
	currentWeight uint8
}

// newWeightedRoundRobin creates a new WeightedRoundRobin instance.
func newWeightedRoundRobin(deployments *[]entity.Deployment) *WeightedRoundRobin {
	w := WeightedRoundRobin{
		deployments:   deployments,
		currentIndex:  0,
		currentWeight: uint8(0),
	}

	return &w
}

// Next implements Scheduler.Next. The Weighted Round Robin load balancing
// distributes load circular and weighted among the service instances.
func (w *WeightedRoundRobin) Next() (*entity.Instance, error) {
	attempts := 0

lookup:
	for attempts < len(*w.deployments) {
		index := w.currentIndex % len(*w.deployments)
		d := (*w.deployments)[index]

		if !d.Instance.Config.IsAttached || !d.Instance.IsAlive {
			w.currentIndex++
			attempts++
			continue lookup
		}

		if d.Node.Config.Weight == w.currentWeight {
			w.currentIndex++
			w.currentWeight = uint8(0)
			attempts++
			continue lookup
		}

		if d.Node.Config.Weight > w.currentWeight {
			w.currentWeight++
			return d.Instance, nil
		}

		attempts++
	}

	return nil, errors.New("no service instance found")
}
