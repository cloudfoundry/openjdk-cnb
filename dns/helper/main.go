/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"net"
	"os"

	"github.com/miekg/dns"
)

func main() {
	security := os.Args[1]

	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code, err := d(security, config)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(code)
	}

	os.Exit(code)
}

func d(security string, config *dns.ClientConfig) (int, error) {
	if !isLinkLocal(config) {
		return 0, nil
	}

	fmt.Println("JVM DNS caching disabled in favor of link-local DNS caching")

	f, err := os.OpenFile(security, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return 1, err
	}
	defer f.Close()

	if _, err := f.WriteString(`
networkaddress.cache.ttl=0
networkaddress.cache.negative.ttl=0
`); err != nil {
		return 1, err
	}

	return 0, nil
}

func isLinkLocal(config *dns.ClientConfig) bool {
	return net.ParseIP(config.Servers[0]).IsLinkLocalUnicast()
}
