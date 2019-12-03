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

// Package api provides an API server for controlling the Dice core.
package api

import (
	"context"
	"github.com/dominikbraun/dice/controller"
	"github.com/go-chi/chi"
	"net/http"
)

// ServerConfig concludes properties that are configurable by the user.
type ServerConfig struct {
	Address string `json:"address"`
	Logfile string `json:"logfile"`
}

// Server is the actual HTTP server exposing a REST API. It will accept
// requests on the specified TCP address and handles these requests using
// the provided controller.Controller instance. The listening port has to
// be secured against remote access.
type Server struct {
	config     ServerConfig
	router     chi.Router
	server     *http.Server
	controller *controller.Controller
}

// NewServer creates a new Server instance and initializes all routes.
func NewServer(config ServerConfig, controller *controller.Controller) *Server {
	s := Server{
		config:     config,
		router:     newRouter(),
		controller: controller,
	}

	s.server = &http.Server{
		Addr:    s.config.Address,
		Handler: s.router,
	}

	return &s
}

// Run makes the API server listen on the specified TCP address and accept
// incoming requests. This function should be called in an extra goroutine
// since Run is a blocking function.
//
// Unlike ListenAndServe from net/http, Run only returns real errors, meaning
// that it won't return an error when shutting down.
func (s *Server) Run() error {
	err := s.server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown attempts a graceful shutdown. Active connections will not be
// interrupted until the context used inside Shutdown expires.
func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}
