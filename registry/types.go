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

// Package registry provides the service registry and the route registry.
//
// While the core package as well as the store package represent the data
// statically and storage-oriented, the registries provide a representation
// required at runtime: In-memory, dynamic and quickly accessible.
package registry

import "github.com/dominikbraun/dice/entity"

// Scheduler represents a load balancing algorithm that manages multiple
// deployments of a service and returns the next instance using `Next`.
type Scheduler interface {
	Next() (*entity.Instance, error)
	UpdateDeployments(deployments []Deployment)
}

// Service is the service representation used by the registries. Compared
// to a entity.Service, it does not hold only meta data like its name but
// also stores all deployments and the associated scheduler.
//
// This representation is important for quick load balancing: The proxy
// asks for the service and gets a Service instance. Using the service's
// scheduler, it can determine the instance for forwarding the request.
type Service struct {
	Entity      *entity.Service
	Deployments []Deployment
	Scheduler   Scheduler
}

// Deployment represents a physical service deployment, simply consisting
// of an instance and the node it has been deployed to. This association
// is used by the scheduler for load balancing.
type Deployment struct {
	Node     *entity.Node
	Instance *entity.Instance
}

// isRemovable checks if a deployment can be removed safely.
func (d Deployment) isRemovable() bool {
	return !d.Node.IsAttached && !d.Instance.IsAttached
}

// equals checks if a deployment is considered equal to another deployment.
func (d Deployment) equals(other Deployment) bool {
	nodeIsEqual := d.Node.ID == other.Node.ID
	instanceIsEqual := d.Instance.ID == other.Instance.ID

	return nodeIsEqual && instanceIsEqual
}
