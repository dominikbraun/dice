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

// Package types provides common types shared across packages.
package types

// NodeCreateOptions combines all user options for creating a new node.
// It serves as a Data Transfer Object for the Dice core.
type NodeCreateOptions struct {
	Weight uint8 `json:"weight"`
	Attach bool  `json:"attach"`
}

// NodeRemoveOptions combines all user options for removing a node.
type NodeRemoveOptions struct {
	Force bool `json:"force"`
}

// NodeInfoOptions combines all user options for printing information
// about a node.
type NodeInfoOptions struct {
	Quiet bool `json:"quiet"`
}

// NodeInfoOptions combines all user options for listing nodes.
type NodeListOptions struct {
	All bool `json:"all"`
}

// ServiceCreateOptions combines all user options for creating a new
// service. It serves as a Data Transfer Object for the Dice core.
type ServiceCreateOptions struct {
	Balancing string `json:"balancing"`
	Enable    bool   `json:"enable"`
}

// ServiceInfoOptions combines all user options for printing information
// about an service.
type ServiceInfoOptions struct {
	Quiet bool `json:"quiet"`
}

// ServiceListOptions combines all user options for listing services.
type ServiceListOptions struct {
	All bool `json:"all"`
}

// ServiceURLOptions combines all user options for setting service URLs.
type ServiceURLOptions struct {
	Delete bool `json:"delete"`
}

// InstanceCreateOptions combines all user options for creating a new
// instance. It serves as a Data Transfer Object for the Dice core.
type InstanceCreateOptions struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Attach  bool   `json:"attach"`
}

// InstanceRemoveOptions combines all user options for removing an
// instance.
type InstanceRemoveOptions struct {
	Force bool `json:"force"`
}

// InstanceInfoOptions combines all user options for printing information
// about an instance.
type InstanceInfoOptions struct {
	Quiet bool `json:"quiet"`
}

// InstanceListOptions combines all user options for listing instances.
type InstanceListOptions struct {
	All bool `json:"all"`
}
