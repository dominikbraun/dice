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

// CreateNode handles a POST request for creating a new node. The request
// body has to contain the node's URL and associated options.
func (c *Controller) CreateNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var nodeCreate types.NodeCreate

		if err := json.NewDecoder(r.Body).Decode(&nodeCreate); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, ErrInvalidFormData)
			return
		}

		if err := c.backend.CreateNode(nodeCreate.Name, nodeCreate.NodeCreateOptions); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusOK, types.Response{Success: true})
	}
}

// AttachNode handles a POST request for attaching an existing node. The
// request URL has to contain a valid node reference.
func (c *Controller) AttachNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeRef := entity.NodeReference(chi.URLParam(r, "ref"))

		if err := c.backend.AttachNode(nodeRef); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
		}

		respond(w, r, http.StatusOK, types.Response{Success: true})
	}
}

// AttachNode handles a POST request for detaching an existing node. The
// request URL has to contain a valid node reference.
func (c *Controller) DetachNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeRef := entity.NodeReference(chi.URLParam(r, "ref"))

		if err := c.backend.DetachNode(nodeRef); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
		}

		respond(w, r, http.StatusOK, types.Response{Success: true})
	}
}

// RemoveNode handles a POST request for removing an existing node. The
// request URL has to contain a valid node reference.
func (c *Controller) RemoveNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeRef := entity.NodeReference(chi.URLParam(r, "ref"))

		var options types.NodeRemoveOptions

		if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := c.backend.RemoveNode(nodeRef, options); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
		}

		respond(w, r, http.StatusOK, types.Response{Success: true})
	}
}

// NodeInfo handles a POST request for retrieving information for a node. The
// request URL has to contain a valid node reference.
func (c *Controller) NodeInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeRef := entity.NodeReference(chi.URLParam(r, "ref"))

		nodeInfo, err := c.backend.NodeInfo(nodeRef)
		if err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusOK, types.Response{Success: true, Data: nodeInfo})
	}
}

// ListNodes handles a POST request for retrieving a list of nodes. The request
// body has to contain valid NodeListOptions.
func (c *Controller) ListNodes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var options types.NodeListOptions

		if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, ErrInvalidFormData)
			return
		}

		nodeList, err := c.backend.ListNodes(options)
		if err != nil {
			respondError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusOK, types.Response{Success: true, Data: nodeList})
	}
}
