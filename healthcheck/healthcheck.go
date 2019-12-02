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

package healthcheck

import (
	"errors"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/registry"
	"net"
	"time"
)

var (
	ErrInvalidDeployments = errors.New("provided deployments are invalid")
)

type Config struct {
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
}

type HealthCheck struct {
	config   Config
	services map[string]registry.Service
	stop     chan bool
}

func New(config Config, services map[string]registry.Service) (*HealthCheck, error) {
	if services == nil {
		return nil, ErrInvalidDeployments
	}

	hc := HealthCheck{
		config:   config,
		services: services,
		stop:     make(chan bool),
	}

	return &hc, nil
}

func (hc *HealthCheck) RunPeriodically() error {
	intervalTick := time.NewTicker(hc.config.Interval)

healthcheck:
	for {
		select {
		case <-intervalTick.C:
			hc.checkServices()
		case <-hc.stop:
			break healthcheck
		}
	}

	return nil
}

func (hc *HealthCheck) checkServices() {
	for _, s := range hc.services {
		if s.Entity.IsEnabled {
			for _, d := range s.Deployments {
				d.Instance.IsAlive = hc.pingInstance(d.Instance)
				// ToDo: If all instances are dead, check if the node is alive
			}
		}
	}
}

func (hc *HealthCheck) pingInstance(instance *entity.Instance) bool {
	address := instance.URL.String()

	conn, err := net.DialTimeout("tcp", address, hc.config.Timeout)
	if err != nil {
		return false
	}

	_ = conn.Close()
	return true
}

func (hc *HealthCheck) Stop() error {
	hc.stop <- true
	return nil
}
