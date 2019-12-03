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
)

func (c *Controller) CreateService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			respond(w, r, http.StatusInternalServerError, ErrInternalServerError.Error())
			return
		}

		name := r.Form.Get("name")

		var options types.ServiceCreateOptions

		if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, ErrInvalidFormData.Error())
			return
		}

		if err := c.backend.CreateService(name, options); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}

		respond(w, r, http.StatusOK, true)
	}
}

func (c *Controller) EnableService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			respond(w, r, http.StatusInternalServerError, ErrInternalServerError.Error())
			return
		}

		serviceRef := entity.ServiceReference(chi.URLParam(r, "ref"))

		if err := c.backend.EnableService(serviceRef); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
		}

		respond(w, r, http.StatusOK, true)
	}
}

func (c *Controller) DisableService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			respond(w, r, http.StatusInternalServerError, ErrInternalServerError.Error())
			return
		}

		serviceRef := entity.ServiceReference(chi.URLParam(r, "ref"))

		if err := c.backend.DisableService(serviceRef); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
		}

		respond(w, r, http.StatusOK, true)
	}
}

func (c *Controller) ServiceInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			respond(w, r, http.StatusInternalServerError, ErrInternalServerError.Error())
			return
		}

		serviceRef := entity.ServiceReference(chi.URLParam(r, "ref"))

		serviceInfo, err := c.backend.ServiceInfo(serviceRef)
		if err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}

		respond(w, r, http.StatusOK, serviceInfo)
	}
}
