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

package entity

import (
	"net/url"
	"time"

	"github.com/sony/sonyflake"
)

// NodeProperty represents a characteristic a node can be identified by.
// However, such a characteristic doesn't have to be a unique identifier.
type NodeProperty uint

const (
	NodeID   NodeProperty = 1
	NodeName NodeProperty = 2
	NodeURL  NodeProperty = 3
)

// NodeConfig concludes all properties that can be configured by the user.
type NodeConfig struct {
	Name       string   `json:"name"`
	URL        *url.URL `json:"url"`
	Weight     uint8    `json:"weight"`
	IsAttached bool     `json:"is_attached"`
}

// Node represents a network node which holds one or more service instances.
// This may be a physical server or a virtual machine for example.
type Node struct {
	ID            string     `json:"id"`
	Config        NodeConfig `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	AttachedSince time.Time  `json:"attached_since"`
	IsAlive       bool       `json:"is_alive"`
}

// NewNode creates a new Node instance and returns a reference to it. Returns
// an error in case the ID cannot be generated.
func NewNode(config NodeConfig) (*Node, error) {
	uuid, err := generateEntityUUID()
	if err != nil {
		return nil, err
	}

	n := Node{
		ID:            uuid,
		Config:        config,
		CreatedAt:     time.Now(),
		AttachedSince: time.Time{},
		IsAlive:       false,
	}

	return &n, nil
}

// generateEntityUUID generates an unique ID in the form `20f8707d6000108`.
// Even though duplicate IDs are very unlikely, there is no guarantee that
// the ID hasn't been generated before. This has to be checked manually.
func generateEntityUUID() (string, error) {
	flakeSettings := sonyflake.Settings{}
	generator := sonyflake.NewSonyflake(flakeSettings)

	uuid, err := generator.NextID()
	if err != nil {
		return "", err
	}

	return string(uuid), nil
}
