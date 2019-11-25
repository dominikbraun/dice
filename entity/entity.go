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

// Entity can be any concrete entity value such as a entity.Node instance.
type Entity interface{}

// Type enumerates all supported entity types such as Node or Service.
type Type uint

const (
	TypeNode     Type = 0
	TypeService  Type = 1
	TypeInstance Type = 2
)

// Property can be any entity concrete entity property instance, e. g. NodeID.
type Property interface{}
