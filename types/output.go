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

// NodeInfoOutput is the output printed by the `node info` command.
type NodeInfoOutput struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	IsAttached bool   `json:"is_attached"`
	IsAlive    bool   `json:"is_alive"`
}

// ServiceInfoOutput is the output printed by the `service info` command.
type ServiceInfoOutput struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	URLs            []string `json:"urls"`
	TargetVersion   string   `json:"target_version"`
	BalancingMethod string   `json:"balancing_method"`
	IsEnabled       bool     `json:"is_enabled"`
}

// InstanceInfoOutput is the output printed by the `instance info` command.
type InstanceInfoOutput struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ServiceID  string `json:"service_id"`
	NodeID     string `json:"node_id"`
	URL        string `json:"url"`
	Version    string `json:"version"`
	IsAttached bool   `json:"is_attached"`
	IsAlive    bool   `json:"is_alive"`
}
