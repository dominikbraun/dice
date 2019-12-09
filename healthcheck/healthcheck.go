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

// Package healthcheck provides types and methods for periodic health checks.
package healthcheck

import (
	"errors"
	"fmt"
	"github.com/dominikbraun/dice/entity"
	"github.com/dominikbraun/dice/registry"
	"net"
	"time"
)

var (
	ErrInvalidDeployments = errors.New("provided deployments are invalid")
)

// Config concludes the user-configurable properties for health checks.
type Config struct {
	Interval time.Duration `json:"interval"`
	// When Timeout expires without response, an instance is considered dead.
	Timeout time.Duration `json:"timeout"`
}

// HealthCheck is a simple health checker that can run checks periodically as
// well as manually. It will ping all instances of a provided service map and
// mark each instance as dead or alive on each check.
type HealthCheck struct {
	config   Config
	services *map[string]registry.Service
	stop     chan bool
}

// New creates a new HealthCheck instance. It will take all service instances
// from a service map into account.
func New(config Config, services *map[string]registry.Service) (*HealthCheck, error) {
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

// RunPeriodically runs periodic health checks that will start every time the
// configured interval expires. This function should run in an own goroutine.
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

// RunManually triggers a manual, single health check. This function should be
// called in an own goroutine as well, since the health check can take a while.
func (hc *HealthCheck) RunManually() error {
	hc.checkServices()
	return nil
}

// checkServices loops over all services and their deployments. Each instance
// will be pinged and marked as dead or alive after the timeout expires.
func (hc *HealthCheck) checkServices() {
	for _, s := range *hc.services {
		if s.Entity.IsEnabled {
			for _, d := range s.Deployments {
				d.Instance.IsAlive = hc.pingInstance(d.Node, d.Instance)
				// ToDo: If all instances are dead, check if the node is alive
			}
		}
	}
}

// pingInstance reads the address from an instance and attempts to establish a
// connection to that address. The dialer will use the configured timeout.
func (hc *HealthCheck) pingInstance(node *entity.Node, instance *entity.Instance) bool {
	address := fmt.Sprintf("%s:%v", node.URL.Hostname(), instance.Port)

	conn, err := net.DialTimeout("tcp", address, hc.config.Timeout)
	if err != nil {
		return false
	}

	_ = conn.Close()
	return true
}

// Stop gracefully stops an health check. Running checks will not be affected.
func (hc *HealthCheck) Stop() error {
	hc.stop <- true
	return nil
}
