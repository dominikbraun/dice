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
	"github.com/go-chi/chi"
	"net/http"
	"os"
)

type ServerConfig struct {
	Address string `json:"address"`
	Logfile string `json:"logfile"`
}

type Server struct {
	config    ServerConfig
	router    chi.Router
	server    *http.Server
	interrupt chan os.Signal
}

func NewServer(config ServerConfig, quit chan os.Signal) *Server {
	s := Server{
		config:    config,
		router:    newRouter(),
		interrupt: quit,
	}

	s.server = &http.Server{
		Addr:    s.config.Address,
		Handler: s.router,
	}

	return &s
}

func (s *Server) Run() chan<- error {
	errors := make(chan error)

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errors <- err
		}
		close(errors)
	}()

	go func() {
		<-s.interrupt
		if err := s.server.Shutdown(context.Background()); err != nil {
			errors <- err
		}
	}()

	return errors
}
