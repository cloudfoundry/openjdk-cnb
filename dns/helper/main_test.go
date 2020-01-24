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
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/miekg/dns"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLinkLocalDNSHelper(t *testing.T) {
	spec.Run(t, "Link-local DNS Helper", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var (
			root     string
			security string
		)

		it.Before(func() {
			root = test.ScratchDir(t, "link-local-dns")

			security = filepath.Join(root, "java.security")
			test.WriteFile(t, security, "test")
		})

		it("does not modify file if not link local", func() {
			config := &dns.ClientConfig{Servers: []string{"1.1.1.1"}}

			g.Expect(d(security, config)).To(gomega.Equal(0))
			g.Expect(security).To(test.HaveContent("test"))
		})

		it("modifies file if link local", func() {
			config := &dns.ClientConfig{Servers: []string{"169.254.0.1"}}

			g.Expect(d(security, config)).To(gomega.Equal(0))
			g.Expect(security).To(test.HaveContent(`test
networkaddress.cache.ttl=0
networkaddress.cache.negative.ttl=0
`))
		})
	}, spec.Report(report.Terminal{}))
}
