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
	"github.com/dominikbraun/dice/types"
	"net/url"
	"time"
)

// InstanceReference is a string that identifies an instance, e. g. an ID.
type InstanceReference string

// Instance represents a service instance. It is an executable associated
// with a service that has been deployed to a node. Any service can have
// multiple instances, allowing redundancy and higher availability.
//
// Like with nodes, attaching an instance to Dice makes it available for
// receiving requests. If the instance has been deployed to a node that is
// currently detached, it won't receive any requests.
type Instance struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ServiceID     string    `json:"service_id"`
	NodeID        string    `json:"node_id"`
	URL           *url.URL  `json:"url"`
	Version       string    `json:"version"`
	IsAttached    bool      `json:"is_attached"`
	IsUpdated     bool      `json:"is_updated"`
	CreatedAt     time.Time `json:"created_at"`
	AttachedSince time.Time `json:"attached_since"`
	IsAlive       bool      `json:"is_alive"`
}

// NewInstance creates a new Instance instance. It doesn't guarantee uniqueness.
func NewInstance(serviceID, nodeID string, url *url.URL, options types.InstanceCreateOptions) (*Instance, error) {
	uuid, err := generateEntityID()
	if err != nil {
		return nil, err
	}

	i := Instance{
		ID:            uuid,
		Name:          options.Name,
		ServiceID:     serviceID,
		NodeID:        nodeID,
		URL:           url,
		Version:       options.Version,
		IsAttached:    options.Attach,
		IsUpdated:     false,
		CreatedAt:     time.Now(),
		AttachedSince: time.Time{},
		IsAlive:       false,
	}

	return &i, nil
}
