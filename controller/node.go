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
	"net/url"
)

// CreateNode handles a POST request for creating a new node. The request
// body has to contain the node's URL and associated options.
func (c *Controller) CreateNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			respond(w, r, http.StatusInternalServerError, ErrInternalServerError.Error())
			return
		}

		nodeURL, err := url.Parse(r.Form.Get("url"))
		if err != nil {
			respond(w, r, http.StatusUnprocessableEntity, ErrInvalidURL.Error())
			return
		}

		var options types.NodeCreateOptions

		if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, ErrInvalidFormData.Error())
			return
		}

		if err := c.backend.CreateNode(nodeURL, options); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}

		respond(w, r, http.StatusOK, true)
	}
}

// AttachNode handles a POST request for attaching an existing node. The
// request URL has to contain a valid node reference.
func (c *Controller) AttachNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeRef := entity.NodeReference(chi.URLParam(r, "ref"))

		if err := c.backend.AttachNode(nodeRef); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
		}

		respond(w, r, http.StatusOK, true)
	}
}

// AttachNode handles a POST request for detaching an existing node. The
// request URL has to contain a valid node reference.
func (c *Controller) DetachNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeRef := entity.NodeReference(chi.URLParam(r, "ref"))

		if err := c.backend.DetachNode(nodeRef); err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
		}

		respond(w, r, http.StatusOK, true)
	}
}

// NodeInfo handles a POST request for retrieving information for a node. The
// request URL has to contain a valid node reference.
func (c *Controller) NodeInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeRef := entity.NodeReference(chi.URLParam(r, "ref"))

		nodeInfo, err := c.backend.NodeInfo(nodeRef)
		if err != nil {
			respond(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}

		respond(w, r, http.StatusOK, nodeInfo)
	}
}
