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

package core

import (
	"github.com/dominikbraun/dice/api"
	"github.com/dominikbraun/dice/config"
	"github.com/dominikbraun/dice/log"
	"github.com/dominikbraun/dice/proxy"
	"github.com/dominikbraun/dice/registry"
	"github.com/dominikbraun/dice/store"
	"os"
	"os/signal"
)

const (
	configName  string = "dice"
	kvStorePath string = "./dice-store"
	logfilePath string = "./dice.log"
)

type Dice struct {
	config    config.Reader
	kvStore   store.EntityStore
	registry  *registry.ServiceRegistry
	interrupt chan os.Signal
	apiServer *api.Server
	proxy     *proxy.Proxy
	logger    log.Logger
}

func NewDice() (*Dice, chan<- os.Signal, error) {
	var d Dice
	var err error

	if d.config, err = config.NewConfig(configName); err != nil {
		return nil, nil, err
	}

	if d.kvStore, err = store.NewKV(kvStorePath); err != nil {
		return nil, nil, err
	}

	d.registry = registry.NewServiceRegistry()

	d.interrupt = make(chan os.Signal)
	signal.Notify(d.interrupt, os.Interrupt)

	d.apiServer = api.NewServer(api.ServerConfig{}, d.interrupt)
	d.proxy = proxy.New(proxy.Config{}, d.registry, d.interrupt)

	logfile, err := os.OpenFile(logfilePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, nil, err
	}

	d.logger = log.NewLogger(logfile, log.InfoLevel)

	return &d, d.interrupt, nil
}
