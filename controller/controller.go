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
	"errors"
	"github.com/dominikbraun/dice/types"
	"github.com/go-chi/render"
	"net/http"
)

var (
	ErrInternalServerError = errors.New("an internal server error occurred")
	ErrInvalidURL          = errors.New("the given URL is not valid")
	ErrInvalidFormData     = errors.New("the provided form data is not valid")
)

// Controller is a REST interface that controls the Dice core. It provides
// HTTP handling methods which will read all required data from the request,
// invoke the core functions and eventually return the core's responses.
type Controller struct {
	backend      Target
	ReloadSignal chan bool
}

// New creates a new Controller instance that uses the provided Target.
func New(backend Target) *Controller {
	c := Controller{
		backend: backend,
	}

	return &c
}

// respond sets an HTTP status code and renders any given response value.
// Note that a return statement is required after calling respond.
func respond(w http.ResponseWriter, r *http.Request, status int, response types.Response) {
	w.WriteHeader(status)
	render.JSON(w, r, response)
}

// respondError does the same as respond, however it takes an error as value
// and creates an appropriate response on its own using that error.
func respondError(w http.ResponseWriter, r *http.Request, status int, err error) {
	response := types.Response{
		Success: false,
		Message: err.Error(),
	}
	respond(w, r, status, response)
}
