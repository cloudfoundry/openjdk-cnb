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

package provider_test

import (
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/openjdk-cnb/provider"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestSecurityProviderConfigurer(t *testing.T) {
	spec.Run(t, "SecurityProviderConfigurer", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		when("Java 8", func() {
			var dep buildpack.Dependency

			it.Before(func() {
				version, err := semver.NewVersion("8.0.212")
				g.Expect(err).NotTo(gomega.HaveOccurred())
				dep = buildpack.Dependency{Version: buildpack.Version{Version: version}}
			})

			it("contributes security provider configurer", func() {
				test.TouchFile(t, f.Build.Buildpack.Root, "bin", "security-provider-configurer")

				b, ok, err := provider.NewSecurityProviderConfigurer(f.Build, dep, "test-security")
				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(ok).To(gomega.BeTrue())

				g.Expect(b.Contribute()).To(gomega.Succeed())

				layer := f.Build.Layers.Layer("security-provider-configurer")
				g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
				g.Expect(filepath.Join(layer.Root, "bin", "security-provider-configurer")).To(gomega.BeARegularFile())
				g.Expect(layer).To(test.HaveProfile("security-provider-classpath", `EXT_DIRS="$JAVA_HOME/lib/ext"

for I in ${SECURITY_PROVIDERS_CLASSPATH//:/$'\n'}; do
  EXT_DIRS="$EXT_DIRS:$(dirname $I)"
done

JAVA_OPTS="$JAVA_OPTS -Djava.ext.dirs=$EXT_DIRS"`))
				g.Expect(layer).To(test.HaveProfile("security-provider-configurer", "security-provider-configurer $JAVA_HOME 8.0.212 test-security $SECURITY_PROVIDERS"))
			})
		})

		when("Java 11", func() {
			var dep buildpack.Dependency

			it.Before(func() {
				version, err := semver.NewVersion("11.0.3")
				g.Expect(err).NotTo(gomega.HaveOccurred())
				dep = buildpack.Dependency{Version: buildpack.Version{Version: version}}
			})

			it("contributes security provider configurer", func() {
				test.TouchFile(t, f.Build.Buildpack.Root, "bin", "security-provider-configurer")

				b, ok, err := provider.NewSecurityProviderConfigurer(f.Build, dep, "test-security")
				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(ok).To(gomega.BeTrue())

				g.Expect(b.Contribute()).To(gomega.Succeed())

				layer := f.Build.Layers.Layer("security-provider-configurer")
				g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
				g.Expect(filepath.Join(layer.Root, "bin", "security-provider-configurer")).To(gomega.BeARegularFile())
				g.Expect(layer).To(test.HaveProfile("security-provider-classpath", "export CLASSPATH=$CLASSPATH:$SECURITY_PROVIDERS_CLASSPATH"))
				g.Expect(layer).To(test.HaveProfile("security-provider-configurer", "security-provider-configurer $JAVA_HOME 11.0.3 test-security $SECURITY_PROVIDERS"))
			})
		})
	}, spec.Report(report.Terminal{}))
}
