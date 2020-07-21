// Copyright 2020 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/brutella/dnssd"
)

func testConns(ctx context.Context, service *dnssd.Service) {
	port := strconv.Itoa(service.Port)
	hostports := []string{
		// DHCP-based DNS in the local network:
		net.JoinHostPort(service.Host, port),

		// Avahi name resolution:
		net.JoinHostPort(service.Host+"."+service.Domain, port),
	}
	for _, ip := range service.IPs {
		// DNSSD-provided IP address:
		hostports = append(hostports, net.JoinHostPort(ip.String(), port))
	}
	var wg sync.WaitGroup
	for _, hostport := range hostports {
		hostport := hostport // copy
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", hostport)
			if err != nil {
				log.Printf("%s: %v", hostport, err)
				return
			}
			defer conn.Close()
			log.Printf("%s: reachable!", hostport)
		}()
	}
	wg.Wait()
}
