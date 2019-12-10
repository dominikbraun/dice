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

package scheduler

import (
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/registry"
	"testing"
)

// TestWeightedRoundRobin_Next tests WeightedRoundRobin.Next. It sets up
// 5 instances deployed to 3 nodes in total and initialized the scheduler
// with these deployments. For each call of Next, the selected instance
// is compared to the expected instance.
//
// The asserted instance IDs are a consequence of the rules described at
// the WeightedRoundRobin implementation.
func TestWeightedRoundRobin_Next(t *testing.T) {
	node1 := &entity.Node{ID: "n1", Weight: 2, IsAttached: true, IsAlive: true}
	node2 := &entity.Node{ID: "n2", Weight: 1, IsAttached: true, IsAlive: true}
	node3 := &entity.Node{ID: "n3", Weight: 1, IsAttached: false, IsAlive: true}

	instance1 := &entity.Instance{ID: "i1", IsAttached: true, IsAlive: true}
	instance2 := &entity.Instance{ID: "i2", IsAttached: false, IsAlive: true}
	instance3 := &entity.Instance{ID: "i3", IsAttached: true, IsAlive: true}
	instance4 := &entity.Instance{ID: "i4", IsAttached: true, IsAlive: false}
	instance5 := &entity.Instance{ID: "i5", IsAttached: true, IsAlive: true}

	deployments := []registry.Deployment{
		{Node: node1, Instance: instance1},
		{Node: node1, Instance: instance2},
		{Node: node2, Instance: instance3},
		{Node: node2, Instance: instance4},
		{Node: node3, Instance: instance5},
	}

	wrr, err := New(&deployments, WeightedRoundRobinBalancing)
	if err != nil {
		t.Error(err)
	}

	assertions := []string{"i1", "i1", "i3", "i5"}

	for run := 0; run < len(assertions); run++ {
		instance, _ := wrr.Next()
		assertedID := assertions[run]

		if instance.ID != assertedID {
			t.Errorf("selected instance %s, expected %s", instance.ID, assertedID)
		}
	}
}
