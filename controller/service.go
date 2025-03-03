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

// CreateService handles a POST request for creating a new service. The
// request body has to contain the service's name and associated options.
func (c *Controller) CreateService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serviceCreate types.ServiceCreate

		if err := json.NewDecoder(r.Body).Decode(&serviceCreate); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, ErrInvalidFormData)
			return
		}

		if err := c.backend.CreateService(serviceCreate.Name, serviceCreate.ServiceCreateOptions); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusOK, types.Response{Success: true})
	}
}

// EnableService handles a POST request for enabling an existing service.
// The request URL has to contain a valid service reference.
func (c *Controller) EnableService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceRef := entity.ServiceReference(chi.URLParam(r, "ref"))

		if err := c.backend.EnableService(serviceRef); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
		}

		respond(w, r, http.StatusOK, types.Response{Success: true})
	}
}

// DisableService handles a POST request for disabling an existing service.
// The request URL has to contain a valid service reference.
func (c *Controller) DisableService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceRef := entity.ServiceReference(chi.URLParam(r, "ref"))

		if err := c.backend.DisableService(serviceRef); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
		}

		respond(w, r, http.StatusOK, types.Response{Success: true})
	}
}

// UpdateService handles a POST request for updating a service. The request
// URL has to contain a valid service reference, the body must provide a
// valid instance of types.ServiceUpdate.
func (c *Controller) UpdateService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceRef := entity.ServiceReference(chi.URLParam(r, "ref"))

		var serviceUpdate types.ServiceUpdate

		if err := json.NewDecoder(r.Body).Decode(&serviceUpdate); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, ErrInvalidFormData)
			return
		}

		if err := c.backend.UpdateService(serviceRef, serviceUpdate.TargetVersion); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
		}

		respond(w, r, http.StatusOK, types.Response{Success: true})
	}
}

// ServiceInfo handles a POST request for retrieving information for a
// service. The request URL has to contain a valid service reference.
func (c *Controller) ServiceInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceRef := entity.ServiceReference(chi.URLParam(r, "ref"))

		serviceInfo, err := c.backend.ServiceInfo(serviceRef)
		if err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusOK, types.Response{Success: true, Data: serviceInfo})
	}
}

// ListServices handles a POST request for retrieving a list of services. The
// request body has to contain valid ServiceListOptions.
func (c *Controller) ListServices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var options types.ServiceListOptions

		if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, ErrInvalidFormData)
			return
		}

		serviceList, err := c.backend.ListServices(options)
		if err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusOK, types.Response{Success: true, Data: serviceList})
	}
}

// SetServiceURL handles a POST request for adding or removing an URL for a
// given service. The request body has to contain a ServiceURL JSON.
func (c *Controller) SetServiceURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceRef := entity.ServiceReference(chi.URLParam(r, "ref"))
		var serviceURL types.ServiceURL

		if err := json.NewDecoder(r.Body).Decode(&serviceURL); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, ErrInvalidFormData)
			return
		}

		err := c.backend.SetServiceURL(serviceRef, serviceURL.URL, serviceURL.ServiceURLOptions)
		if err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusOK, types.Response{Success: true})
	}
}
