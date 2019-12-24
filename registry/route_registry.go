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

// Package registry provides the service registry and the route registry.
//
// While the core package as well as the store package represent the data
// statically and storage-oriented, the registries provide a representation
// required at runtime: In-memory, dynamic and quickly accessible.
package registry

import "errors"

// ServiceRoute is a host with an optional route that serves as a HTTP
// request target. It is an unique identifier for services.
type ServiceRoute string

var (
	ErrUnregisteredRoute      = errors.New("route is not registered")
	ErrRouteAlreadyRegistered = errors.New("route is already registered")
)

// RouteRegistry is the global registry for service routes. It manages a
// simple mapping between a service route and a corresponding service ID.
type RouteRegistry struct {
	routes map[ServiceRoute]string
}

// NewRouteRegistry creates a new, ready to go RouteRegistry instance.
func NewRouteRegistry() *RouteRegistry {
	rr := RouteRegistry{
		routes: make(map[ServiceRoute]string),
	}

	return &rr
}

// RegisterRoute registers a new route and maps it against a service ID.
// Returns an error if it already exists, unless force is set to `true`.
func (rr *RouteRegistry) RegisterRoute(route string, serviceID string, force bool) error {
	if _, exists := rr.routes[ServiceRoute(route)]; exists {
		if !force {
			return ErrRouteAlreadyRegistered
		}
	}

	rr.routes[ServiceRoute(route)] = serviceID

	return nil
}

// UnregisterRoute removes a route from the registry. Returns an error if
// the route doesn't exist.
func (rr *RouteRegistry) UnregisterRoute(route string) error {
	if _, exists := rr.routes[ServiceRoute(route)]; !exists {
		return ErrUnregisteredRoute
	}
	delete(rr.routes, ServiceRoute(route))

	return nil
}

// LookupServiceID looks up a service ID that is associated with the given
// route. The second return value indicates whether the service ID could
// be found or not.
func (rr *RouteRegistry) LookupServiceID(route string) (string, bool) {
	if serviceID, exists := rr.routes[ServiceRoute(route)]; exists {
		return serviceID, true
	}

	return "", false
}

// IsRegistered checks and returns if a given route is registered. Note
// that there's a difference between `example.com` and `example.com/`.
func (rr *RouteRegistry) IsRegistered(route string) bool {
	_, exists := rr.routes[ServiceRoute(route)]
	return exists
}
