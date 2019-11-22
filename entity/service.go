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

// ServiceProperty represents a characteristic a service can be identified by.
// However, such a characteristic doesn't have to be a unique identifier.
type ServiceProperty uint

const (
	ServiceID   ServiceProperty = 0
	ServiceName ServiceProperty = 1
)

// ServiceConfig concludes all properties that can be configured by the user.
type ServiceConfig struct {
	Name            string   `json:"name"`
	Hosts           []string `json:"hosts"`
	Weight          uint8    `json:"weight"`
	TargetVersion   string   `json:"target_version"`
	BalancingMethod string   `json:"balancing_method"`
	IsActive        bool     `json:"is_active"`
}

// Service represents a single service of an application or even the application
// itself. For example, the  authentication service of a bookstore app - or even
// the whole bookstore app - can be represented as a service.
type Service struct {
	ID     string        `json:"id"`
	Config ServiceConfig `json:"config"`
}

// NewService creates a new Service instance and returns a reference to it. Returns
// an error in case the ID cannot be generated.
func NewService(config ServiceConfig) (*Service, error) {
	uuid, err := generateEntityUUID()
	if err != nil {
		return nil, err
	}

	s := Service{
		ID:     uuid,
		Config: config,
	}

	return &s, nil
}
