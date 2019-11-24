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

// Package core provides the Dice load balancer.
package core

import (
	"github.com/dominikbraun/dice/registry"
	"github.com/dominikbraun/dice/server"
	"github.com/dominikbraun/dice/storage"
)

// Dice represents the Dice load balancer and groups all components.
type Dice struct {
	kvStore     storage.EntityStorage
	memory      storage.EntityStorage
	registry    registry.ServiceRegistry
	apiServer   server.APIServer
	proxyServer server.ProxyServer
}

// NewDice creates a new Dice instance and initializes all components.
func NewDice() *Dice {
	var d Dice

	// Initialize components...

	return &d
}
