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

package registry

import "errors"

type ServiceRoute string

var (
	ErrUnregisteredRoute      = errors.New("route is not registered")
	ErrRouteAlreadyRegistered = errors.New("route is already registered")
)

type RouteRegistry struct {
	routes map[ServiceRoute]string
}

func (rr *RouteRegistry) RegisterRoute(route string, serviceID string, force bool) error {
	if _, exists := rr.routes[ServiceRoute(route)]; exists {
		if !force {
			return ErrRouteAlreadyRegistered
		}
	}

	rr.routes[ServiceRoute(route)] = serviceID

	return nil
}

func (rr *RouteRegistry) UnregisterRoute(route string) error {
	if _, exists := rr.routes[ServiceRoute(route)]; !exists {
		return ErrUnregisteredRoute
	}
	delete(rr.routes, ServiceRoute(route))

	return nil
}

func (rr *RouteRegistry) FindServiceID(route string) string {
	if serviceID, exists := rr.routes[ServiceRoute(route)]; exists {
		return serviceID
	}

	return ""
}

func (rr *RouteRegistry) IsRegistered(route string) bool {
	_, exists := rr.routes[ServiceRoute(route)]
	return exists
}
