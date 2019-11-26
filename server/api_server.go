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

package server

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"net/http"
	"os"
)

type APIServerConfig struct {
	Address string `json:"address"`
	Logfile string `json:"logfile"`
}

type APIServer struct {
	config    APIServerConfig
	router    chi.Router
	server    *http.Server
	interrupt chan os.Signal
}

func NewAPIServer(config APIServerConfig, quit chan os.Signal) *APIServer {
	as := APIServer{
		config:    config,
		router:    buildRouter(),
		interrupt: quit,
	}

	as.server = &http.Server{
		Addr:    as.config.Address,
		Handler: as.router,
	}

	return &as
}

func (as *APIServer) Run() chan<- error {
	errors := make(chan error)

	go func() {
		err := as.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errors <- err
		}
		close(errors)
	}()

	go func() {
		<-as.interrupt
		if err := as.server.Shutdown(context.Background()); err != nil {
			errors <- err
		}
	}()

	return errors
}

func buildRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		render.SetContentType(render.ContentTypeJSON),
	)

	return r
}
