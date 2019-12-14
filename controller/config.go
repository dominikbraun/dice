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

// Package controller provides methods for handling REST requests.
package controller

import (
	"github.com/dominikbraun/dice/types"
	"net/http"
)

// ReloadConfig handles a POST request for reloading the configuration.
func (c *Controller) ReloadConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.ReloadSignal <- true
		respond(w, r, http.StatusOK, types.Response{Success: true})
	}
}
