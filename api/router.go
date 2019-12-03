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
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func newRouter() chi.Router {
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

func (s *Server) mountRoutes() {
	r := chi.NewRouter()

	r.Route("/nodes", func(r chi.Router) {
		r.Post("/create", s.controller.CreateNode())

		r.Route("/{ref}", func(r chi.Router) {
			r.Post("/attach", s.controller.AttachNode())
			r.Post("/detach", s.controller.DetachNode())
			r.Post("/info", s.controller.NodeInfo())
		})
	})

	r.Route("/services", func(r chi.Router) {
		r.Post("/create", s.controller.CreateService())

		r.Route("/{ref}", func(r chi.Router) {
			r.Post("/enable", s.controller.EnableService())
			r.Post("/disable", s.controller.DisableService())
			r.Post("/info", s.controller.ServiceInfo())
		})
	})

	r.Route("/instances", func(r chi.Router) {
		r.Post("/create", s.controller.CreateInstance())

		r.Route("/{ref}", func(r chi.Router) {
			r.Post("/attach", s.controller.AttachInstance())
			r.Post("/detach", s.controller.DetachInstance())
			r.Post("/info", s.controller.InstanceInfo())
		})
	})

	s.router.Mount("/v1", r)
}
