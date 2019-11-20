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
)

// InstanceProperty represents a characteristic an instance can be identified by.
// However, such a characteristic doesn't have to be a unique identifier.
type InstanceProperty uint

const (
	InstanceID        InstanceProperty = 0
	InstanceName      InstanceProperty = 1
	InstanceServiceID InstanceProperty = 2
	InstanceURL       InstanceProperty = 3
)

// InstanceConfig concludes all properties that can be configured by the user.
type InstanceConfig struct {
	Name       string   `json:"name"`
	ServiceID  string   `json:"service_id"`
	URL        *url.URL `json:"url"`
	Version    string   `json:"version"`
	IsAttached bool     `json:"is_attached"`
	IsUpdated  bool     `json:"is_updated"`
}

// Instance represents a service instance. While a service merely is a local
// unit in Dice, an instance is an executable instance of such a service.
//
// While an instance can only belong to one service, a service typically has
// multiple instances if the app is non-monolithic. Service instances can
// be distributed over multiple nodes since this is a common scenario.
type Instance struct {
	ID            string         `json:"id"`
	Config        InstanceConfig `json:"config"`
	CreatedAt     time.Time      `json:"created_at"`
	AttachedSince time.Time      `json:"attached_since"`
	IsAlive       bool           `json:"is_alive"`
}

// NewInstance creates a new Instance instance and returns a reference to it. Returns
// an error in case the ID cannot be generated.
func NewInstance(config InstanceConfig) (*Instance, error) {
	uuid, err := generateEntityUUID()
	if err != nil {
		return nil, err
	}

	i := Instance{
		ID:            uuid,
		Config:        config,
		CreatedAt:     time.Now(),
		AttachedSince: time.Time{},
		IsAlive:       false,
	}

	return &i, nil
}
