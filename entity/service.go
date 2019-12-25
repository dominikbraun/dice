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
	"fmt"
	"github.com/dominikbraun/dice/types"
)

// ServiceReference is a string that identifies a service, e. g. an ID.
type ServiceReference string

// Service represents an application or webservice. A Service itself is not
// a running application. Instead, the running executables are represented
// by service instances (see entity.Instance).
//
// Each service is available under multiple URLs like api.example.com and
// example.com/api. Also, the load balancing algorithm is configurable for
// each service. If a service is disabled, requests will run into HTTP 503.
type Service struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	URLs            []string `json:"urls"`
	TargetVersion   string   `json:"target_version"`
	BalancingMethod string   `json:"balancing_method"`
	IsEnabled       bool     `json:"is_enabled"`
}

// NewService creates a new Service instance. It doesn't guarantee uniqueness.
func NewService(name string, options types.ServiceCreateOptions) (*Service, error) {
	uuid, err := generateEntityID()
	if err != nil {
		return nil, err
	}

	s := Service{
		ID:              uuid,
		Name:            name,
		URLs:            make([]string, 0),
		TargetVersion:   "",
		BalancingMethod: options.Balancing,
		IsEnabled:       options.Enable,
	}

	return &s, nil
}

// AddURL adds a public URL to a service.
func (s *Service) AddURL(url string) error {
	index := s.indexOfURL(url)

	if index != -1 {
		return fmt.Errorf("URL '%s' is already registered", url)
	}

	s.URLs = append(s.URLs, url)
	return nil
}

// RemoveURL removes a public URL from a service.
func (s *Service) RemoveURL(url string) error {
	index := s.indexOfURL(url)

	if index == -1 {
		return fmt.Errorf("URL '%s' is not registered", url)
	}

	urls := s.URLs
	urls[index] = urls[len(urls)-1]
	s.URLs = urls[:len(urls)-1]

	return nil
}

// indexOfURL determines the index of a given URL in the `URLs` field.
func (s *Service) indexOfURL(url string) int {
	index := -1

	for i, u := range s.URLs {
		if u == url {
			index = i
		}
	}

	return index
}
