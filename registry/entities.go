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

package registry

import (
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/scheduler"
)

type Service struct {
	entity      *entity.Service
	deployments []Deployment
	scheduler   scheduler.Scheduler
}

func (s Service) isRemovable() bool {
	for _, d := range s.deployments {
		if !d.isRemovable() {
			return false
		}
	}

	return true
}

type Deployment struct {
	Node     *entity.Node
	Instance *entity.Instance
	IsValid  bool
}

func (d Deployment) isRemovable() bool {
	if d.Node.IsAttached && d.Instance.IsAttached {
		return false
	}

	return true
}

func (d Deployment) equals(other Deployment) bool {
	nodeIsEqual := d.Node.ID == other.Node.ID
	instanceIsEqual := d.Instance.ID == other.Instance.ID

	return nodeIsEqual && instanceIsEqual
}
