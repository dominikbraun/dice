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

// Package controller provides methods for handling REST requests.
package controller

import (
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/types"
	"net/url"
)

// Target concludes all *Target interfaces. Any Target implementation is
// a wholesome implementation of the Dice core.
type Target interface {
	NodeTarget
	ServiceTarget
	InstanceTarget
}

// NodeTarget prescribes methods for backends working with nodes.
type NodeTarget interface {
	CreateNode(url *url.URL, options types.NodeCreateOptions) error
	AttachNode(nodeRef entity.NodeReference) error
	DetachNode(nodeRef entity.NodeReference) error
	NodeInfo(nodeRef entity.NodeReference) (types.NodeInfoOutput, error)
	ListNodes(options types.NodeListOptions) ([]types.NodeInfoOutput, error)
}

// ServiceTarget prescribes methods for backends working with services.
type ServiceTarget interface {
	CreateService(name string, options types.ServiceCreateOptions) error
	EnableService(serviceRef entity.ServiceReference) error
	DisableService(serviceRef entity.ServiceReference) error
	ServiceInfo(serviceRef entity.ServiceReference) (types.ServiceInfoOutput, error)
}

// InstanceTarget prescribes methods for backends working with instances.
type InstanceTarget interface {
	CreateInstance(serviceRef entity.ServiceReference, nodeRef entity.NodeReference, url string, options types.InstanceCreateOptions) error
	AttachInstance(instanceRef entity.InstanceReference) error
	DetachInstance(instanceRef entity.InstanceReference) error
	InstanceInfo(instanceRef entity.InstanceReference) (types.InstanceInfoOutput, error)
}
