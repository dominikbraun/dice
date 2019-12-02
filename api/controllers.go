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

import "net/http"

type Controller interface {
	NodeController
	ServiceController
	InstanceController
}

type NodeController interface {
	CreateNode() http.HandlerFunc
	AttachNode() http.HandlerFunc
	DetachNode() http.HandlerFunc
	NodeInfo() http.HandlerFunc
}

type ServiceController interface {
	CreateService() http.HandlerFunc
	EnableService() http.HandlerFunc
	DisableService() http.HandlerFunc
	ServiceInfo() http.HandlerFunc
}

type InstanceController interface {
	CreateInstance() http.HandlerFunc
	AttachInstance() http.HandlerFunc
	DetachInstance() http.HandlerFunc
	InstanceInfo() http.HandlerFunc
}
