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

package types

import (
	"github.com/dominikbraun/dice/entity"
	"net/url"
)

type NodeInfoOutput struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	URL        *url.URL `json:"url"`
	IsAttached bool     `json:"is_attached"`
	IsAlive    bool     `json:"is_alive"`
}

func (ni NodeInfoOutput) Populate(node *entity.Node) {
	ni.ID = node.ID
	ni.Name = node.Name
	ni.URL = node.URL
	ni.IsAttached = node.IsAttached
	ni.IsAttached = node.IsAlive
}
