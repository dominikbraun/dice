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

// Package core provides the Dice load balancer and its methods.
package core

import (
	"github.com/dominikbraun/dice/entity"
	"regexp"
)

// urlSafe specifies a regular expression for a valid URL. It only allows
// characters that are URL-safe according to RFC 3986.
var urlSafe = regexp.MustCompile("^[a-zA-Z0-9_-]*$")

// validateNode checks all node properties and determines if they're valid.
// It does not check whether the node does already exist or not.
func validateNode(node *entity.Node) (bool, string) {
	if !urlSafe.MatchString(node.ID) {
		return false, "ID must only contain _ and - as special characters"
	}

	if !urlSafe.MatchString(node.Name) {
		return false, "Name must only contain _ and - as special characters"
	}

	return true, ""
}

// validateService checks all service properties and determines if they're
// valid. It does not check whether the service does already exist or not.
func validateService(service *entity.Service) (bool, string) {
	if !urlSafe.MatchString(service.ID) {
		return false, "ID must only contain _ and - as special characters"
	}

	if service.Name == "" {
		return false, "Name must not be empty"
	}

	if !urlSafe.MatchString(service.Name) {
		return false, "Name must only contain _ and - as special characters"
	}

	return true, ""
}

// validateInstance checks all instance properties and determines if they're
// valid. It does not check whether the instance does already exist or not.
func validateInstance(instance *entity.Instance) (bool, string) {
	if !urlSafe.MatchString(instance.ID) {
		return false, "ID must only contain _ and - as special characters"
	}

	if !urlSafe.MatchString(instance.Name) {
		return false, "Name must only contain _ and - as special characters"
	}

	return true, ""
}
