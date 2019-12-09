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

// Package types provides common types shared across packages.
package types

// NodeCreate is a type exclusively used for the REST API. It holds all
// information required to create a new node.
//
// The Dice core asks for an URL and some options when creating a new node.
// It takes the URL and the options as separate parameters. However, for the
// API use case, it is significantly simpler to read the request's JSON into
// a single struct. This is NodeCreate's purpose: It combines all necessary
// information and gets populated with the JSON from the request body.
type NodeCreate struct {
	URL string `json:"url"`
	NodeCreateOptions
}

// ServiceCreate is a type exclusively used for the REST API. It holds all
// information required to create a new service.
//
// For further information, see the documentation for NodeCreate.
type ServiceCreate struct {
	Name string `json:"name"`
	ServiceCreateOptions
}

// InstanceCreate is a type exclusively used for the REST API. It holds all
// information required to create a new instance.
//
// For further information, see the documentation for NodeCreate.
type InstanceCreate struct {
	ServiceRef string `json:"service_ref"`
	NodeRef    string `json:"node_ref"`
	Port       uint16 `json:"port"`
	InstanceCreateOptions
}
