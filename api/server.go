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

package api

import (
	"context"
	"github.com/dominikbraun/dice/controller"
	"github.com/go-chi/chi"
	"net/http"
)

type ServerConfig struct {
	Address string `json:"address"`
	Logfile string `json:"logfile"`
}

type Server struct {
	config     ServerConfig
	router     chi.Router
	server     *http.Server
	controller *controller.Controller
}

func NewServer(config ServerConfig, backend controller.Target) *Server {
	s := Server{
		config: config,
		router: newRouter(),
	}

	s.server = &http.Server{
		Addr:    s.config.Address,
		Handler: s.router,
	}

	s.controller = controller.New(backend)

	return &s
}

func (s *Server) Run() error {
	err := s.server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}
