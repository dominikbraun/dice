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
	Name   string `json:"name"`
	Weight uint8  `json:"weight"`
	Attach bool   `json:"attach"`
}

// NodeInfoOptions combines all user options for printing information
// about a node.
type NodeInfoOptions struct {
	Quiet bool `json:"quiet"`
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

// InstanceCreateOptions combines all user options for creating a new
// instance. It serves as a Data Transfer Object for the Dice core.
type InstanceCreateOptions struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Attach  bool   `json:"attach"`
}

// InstanceInfoOptions combines all user options for printing information
// about an instance.
type InstanceInfoOptions struct {
	Quiet bool `json:"quiet"`
}
