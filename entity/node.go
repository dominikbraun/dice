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
	"github.com/sony/sonyflake"
	"net/url"
	"time"
)

type NodeConfig struct {
	Name       string   `json:"name"`
	URL        *url.URL `json:"url"`
	Weight     uint8    `json:"weight"`
	IsAttached bool     `json:"is_attached"`
}

type Node struct {
	ID            string     `json:"id"`
	Config        NodeConfig `json:"config"`
	CreatedAt     time.Time  `json:"created_at"`
	AttachedSince time.Time  `json:"attached_since"`
	IsAlive       bool       `json:"is_alive"`
}

func NewNode(config NodeConfig) (*Node, error) {
	uuid, err := generateEntityID()
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

func generateEntityID() (string, error) {
	flakeSettings := sonyflake.Settings{}
	generator := sonyflake.NewSonyflake(flakeSettings)

	uuid, err := generator.NextID()
	if err != nil {
		return "", err
	}

	return string(uuid), nil
}
