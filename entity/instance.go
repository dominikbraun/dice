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

type InstanceConfig struct {
	Name       string   `json:"name"`
	ServiceID  string   `json:"service_id"`
	NodeID     string   `json:"node_id"`
	URL        *url.URL `json:"url"`
	Version    string   `json:"version"`
	IsAttached bool     `json:"is_attached"`
	IsUpdated  bool     `json:"is_updated"`
}

type Instance struct {
	ID            string         `json:"id"`
	Config        InstanceConfig `json:"config"`
	CreatedAt     time.Time      `json:"created_at"`
	AttachedSince time.Time      `json:"attached_since"`
	IsAlive       bool           `json:"is_alive"`
}

func NewInstance(config InstanceConfig) (*Instance, error) {
	uuid, err := generateEntityID()
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
