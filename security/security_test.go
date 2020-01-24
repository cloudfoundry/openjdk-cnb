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

package security_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/cloudfoundry/openjdk-cnb/security"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestSecurity(t *testing.T) {
	spec.Run(t, "Java Security Properties", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("contributes java security properties", func() {
			f := test.NewBuildFactory(t)

			s := security.NewSecurity(f.Build)

			g.Expect(s.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("java-security-properties")
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))

			destination := filepath.Join(layer.Root, "java.security")
			g.Expect(destination).To(gomega.BeARegularFile())
			g.Expect(layer).To(test.HaveAppendLaunchEnvironment("JAVA_OPTS", " -Djava.security.properties=%s", destination))
		})
	}, spec.Report(report.Terminal{}))
}
