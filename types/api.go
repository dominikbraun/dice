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
//
// The other following types like ServiceCreate serve the exact same purpose.
type NodeCreate struct {
	Name string `json:"name"`
	NodeCreateOptions
}

// ServiceCreate is a type exclusively used for the REST API. It holds all
// information required to create a new service.
//
// For further information about its usage, see the docs for NodeCreate.
type ServiceCreate struct {
	Name string `json:"name"`
	ServiceCreateOptions
}

// ServiceUpdate is a type exclusively used for the REST API. It holds all
// information required to update a service.
//
// For further information about its usage, see the docs for NodeCreate.
type ServiceUpdate struct {
	TargetVersion string `json:"target_version"`
}

// ServiceURL is a type exclusively used for the REST API. It holds all
// information required to set an URL for a service.
//
// For further information about its usage, see the docs for NodeCreate.
type ServiceURL struct {
	URL string `json:"url"`
	ServiceURLOptions
}

// InstanceCreate is a type exclusively used for the REST API. It holds all
// information required to create a new instance.
//
// For further information about its usage, see the docs for NodeCreate.
type InstanceCreate struct {
	ServiceRef string `json:"service_ref"`
	NodeRef    string `json:"node_ref"`
	URL        string `json:"url"`
	InstanceCreateOptions
}

// Response represents an API response that will be returned to the client.
//
// All *Response types wrap this basic response and a specific *Output type,
// forming an API response for a specific command.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NodeInfoResponse is an API response that carries a NodeInfoOutput.
type NodeInfoResponse struct {
	Response
	Data NodeInfoOutput `json:"data"`
}

// NodeListResponse is an API response that carries a list of nodes. At the
// moment, this is a list of NodeInfoOutputs as returned by the Dice core.
type NodeListResponse struct {
	Response
	Data []NodeInfoOutput `json:"data"`
}

// ServiceInfoResponse carrying a ServiceInfoOutput.
type ServiceInfoResponse struct {
	Response
	Data ServiceInfoOutput `json:"data"`
}

// ServiceListResponse is an API response that carries a list of services.
// At the moment, this is a list of ServiceInfoOutputs as returned by the
// Dice core.
type ServiceListResponse struct {
	Response
	Data []ServiceInfoOutput `json:"data"`
}

// InstanceInfoResponse carrying a InstanceInfoOutput.
type InstanceInfoResponse struct {
	Response
	Data InstanceInfoOutput `json:"data"`
}

// InstanceListResponse is an API response that carries a list of instances.
// At the moment, this is a list of InstanceInfoOutputs as returned by the
// Dice core.
type InstanceListResponse struct {
	Response
	Data []InstanceInfoOutput `json:"data"`
}
