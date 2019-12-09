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
	"encoding/json"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/types"
	"github.com/go-chi/chi"
	"net/http"
)

// CreateInstance handles a POST request for creating a new instance. The
// request body has to contain a valid service reference, node reference
// and instance URL as well as associated options.
func (c *Controller) CreateInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var instanceCreate types.InstanceCreate

		if err := json.NewDecoder(r.Body).Decode(&instanceCreate); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, ErrInvalidFormData.Error())
			return
		}

		serviceRef := entity.ServiceReference(instanceCreate.ServiceRef)
		nodeRef := entity.NodeReference(instanceCreate.NodeRef)

		if err := c.backend.CreateInstance(serviceRef, nodeRef, instanceCreate.Port, instanceCreate.InstanceCreateOptions); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}

		respond(w, r, http.StatusOK, true)
	}
}

// AttachInstance handles a POST request for attaching an existing instance.
// The request URL has to contain a valid instance reference.
func (c *Controller) AttachInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instanceRef := entity.InstanceReference(chi.URLParam(r, "ref"))

		if err := c.backend.AttachInstance(instanceRef); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
		}

		respond(w, r, http.StatusOK, true)
	}
}

// DetachInstance handles a POST request for detaching an existing instance.
// The request URL has to contain a valid instance reference.
func (c *Controller) DetachInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instanceRef := entity.InstanceReference(chi.URLParam(r, "ref"))

		if err := c.backend.DetachInstance(instanceRef); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
		}

		respond(w, r, http.StatusOK, true)
	}
}

// InstanceInfo handles a POST request for retrieving information for an
// instance. The request URL has to contain a valid node reference.
func (c *Controller) InstanceInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instanceRef := entity.InstanceReference(chi.URLParam(r, "ref"))

		instanceInfo, err := c.backend.InstanceInfo(instanceRef)
		if err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}

		respond(w, r, http.StatusOK, instanceInfo)
	}
}
