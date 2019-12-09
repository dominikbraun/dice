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

// Package entity provides domain entities and their factory functions.
package entity

import (
	"crypto/rand"
	"fmt"
	"github.com/dominikbraun/dice/types"
	"net/url"
	"time"
)

// NodeReference is a string that identifies a node, e. g. an ID or name.
type NodeReference string

// Node represents a network node that one or more applications run on,
// for example a physical server, virtual machine or even a container.
//
// Each node has a weight depicting the node's physical computing power.
// The heavier a node is, the more requests it receives from Dice. Each
// node can be attached to Dice, making it available for these requests.
type Node struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	URL           *url.URL  `json:"url"`
	Weight        uint8     `json:"weight"`
	IsAttached    bool      `json:"is_attached"`
	CreatedAt     time.Time `json:"created_at"`
	AttachedSince time.Time `json:"attached_since"`
	IsAlive       bool      `json:"is_alive"`
}

// NewNode creates a new Node instance. It doesn't guarantee uniqueness.
func NewNode(url *url.URL, options types.NodeCreateOptions) (*Node, error) {
	uuid, err := generateEntityID()
	if err != nil {
		return nil, err
	}

	n := Node{
		ID:            uuid,
		Name:          options.Name,
		URL:           url,
		Weight:        options.Weight,
		IsAttached:    options.Attach,
		CreatedAt:     time.Now(),
		AttachedSince: time.Time{},
		IsAlive:       false,
	}

	return &n, nil
}

// generateEntityID generates a random, time-based ID. Even though an ID
// collision is unlikely, you have to check if the ID really is unique.
func generateEntityID() (string, error) {
	b := make([]byte, 8)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x", b)

	return uuid, nil
}
