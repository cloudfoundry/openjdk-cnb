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

	"github.com/Masterminds/semver"
	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestSecurityProviderConfigurerHelper(t *testing.T) {
	spec.Run(t, "Security Provider Configurer Helper", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var (
			out  string
			root string
		)

		it.Before(func() {
			root = test.ScratchDir(t, "security-provider-configurer")

			out = filepath.Join(root, "java.security")
			test.WriteFile(t, out, "test")
		})

		when("Java 8", func() {
			var (
				in      string
				version *semver.Version
			)

			it.Before(func() {
				in = filepath.Join(root, "lib", "security", "java.security")
				test.WriteFile(t, in, `security.provider.1=ALPHA
security.provider.2=BRAVO
security.provider.3=CHARLIE
`)

				v, err := semver.NewVersion("1.8.192")
				g.Expect(err).NotTo(gomega.HaveOccurred())
				version = v
			})

			it("does not modify file if no additions", func() {
				g.Expect(p(out, root, version, []string{})).To(gomega.Equal(0))
				g.Expect(out).To(test.HaveContent("test"))
			})

			it("modifies files if additions", func() {
				g.Expect(p(out, root, version, []string{ "2|DELTA", "ECHO"})).To(gomega.Equal(0))
				g.Expect(out).To(test.HaveContent(`test
security.provider.1=ALPHA
security.provider.2=DELTA
security.provider.3=BRAVO
security.provider.4=CHARLIE
security.provider.5=ECHO
`))
			})
		})

		when("Java 11", func() {
			var (
				in      string
				version *semver.Version
			)

			it.Before(func() {
				in = filepath.Join(root, "conf", "security", "java.security")
				test.WriteFile(t, in, `security.provider.1=ALPHA
security.provider.2=BRAVO
security.provider.3=CHARLIE
`)

				v, err := semver.NewVersion("11.0.3")
				g.Expect(err).NotTo(gomega.HaveOccurred())
				version = v
			})

			it("does not modify file if no additions", func() {
				g.Expect(p(out, root, version, []string{})).To(gomega.Equal(0))
				g.Expect(out).To(test.HaveContent("test"))
			})

			it("modifies files if additions", func() {
				g.Expect(p(out, root, version, []string{ "2|DELTA", "ECHO"})).To(gomega.Equal(0))
				g.Expect(out).To(test.HaveContent(`test
security.provider.1=ALPHA
security.provider.2=DELTA
security.provider.3=BRAVO
security.provider.4=CHARLIE
security.provider.5=ECHO
`))
			})

		})
	}, spec.Report(report.Terminal{}))
}
