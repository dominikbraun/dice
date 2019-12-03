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

package controller

import (
	"encoding/json"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/types"
	"github.com/go-chi/chi"
	"net/http"
	"net/url"
)

func (c *Controller) CreateInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			respond(w, r, http.StatusInternalServerError, ErrInternalServerError.Error())
			return
		}

		serviceRef := entity.ServiceReference(r.Form.Get("service_ref"))
		nodeRef := entity.NodeReference(r.Form.Get("node_ref"))

		instanceURL, err := url.Parse(r.Form.Get("url"))
		if err != nil {
			respond(w, r, http.StatusUnprocessableEntity, ErrInvalidURL.Error())
			return
		}

		var options types.InstanceCreateOptions

		if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, ErrInvalidFormData.Error())
			return
		}

		if err := c.backend.CreateInstance(serviceRef, nodeRef, instanceURL, options); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}

		respond(w, r, http.StatusOK, true)
	}
}

func (c *Controller) AttachInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			respond(w, r, http.StatusInternalServerError, ErrInternalServerError.Error())
			return
		}

		instanceRef := entity.InstanceReference(chi.URLParam(r, "ref"))

		if err := c.backend.AttachInstance(instanceRef); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
		}

		respond(w, r, http.StatusOK, true)
	}
}

func (c *Controller) DetachInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			respond(w, r, http.StatusInternalServerError, ErrInternalServerError.Error())
			return
		}

		instanceRef := entity.InstanceReference(chi.URLParam(r, "ref"))

		if err := c.backend.DetachInstance(instanceRef); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
		}

		respond(w, r, http.StatusOK, true)
	}
}

func (c *Controller) InstanceInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			respond(w, r, http.StatusInternalServerError, ErrInternalServerError.Error())
			return
		}

		instanceRef := entity.InstanceReference(chi.URLParam(r, "ref"))

		instanceInfo, err := c.backend.InstanceInfo(instanceRef)
		if err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}

		respond(w, r, http.StatusOK, instanceInfo)
	}
}
