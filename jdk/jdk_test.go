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

package jdk_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/openjdk-cnb/jdk"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestJDK(t *testing.T) {
	spec.Run(t, "JDK", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("returns true if build plan exists", func() {
			f.AddDependency(jdk.Dependency, filepath.Join("testdata", "stub-openjdk-jdk.tar.gz"))
			f.AddPlan(buildpackplan.Plan{Name: jdk.Dependency})

			_, ok, err := jdk.NewJDK(f.Build)
			g.Expect(ok).To(gomega.BeTrue())
			g.Expect(err).NotTo(gomega.HaveOccurred())
		})

		it("returns false if build plan does not exist", func() {
			_, ok, err := jdk.NewJDK(f.Build)
			g.Expect(ok).To(gomega.BeFalse())
			g.Expect(err).NotTo(gomega.HaveOccurred())
		})

		it("contributes JDK", func() {
			f.AddDependency(jdk.Dependency, filepath.Join("testdata", "stub-openjdk-jdk.tar.gz"))
			f.AddPlan(buildpackplan.Plan{Name: jdk.Dependency})

			j, _, err := jdk.NewJDK(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(j.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("openjdk-jdk")
			g.Expect(layer).To(test.HaveLayerMetadata(true, true, false))
			g.Expect(filepath.Join(layer.Root, "fixture-marker")).To(gomega.BeARegularFile())
			g.Expect(layer).To(test.HaveOverrideBuildEnvironment("JAVA_HOME", layer.Root))
			g.Expect(layer).To(test.HaveOverrideBuildEnvironment("JDK_HOME", layer.Root))
		})
	}, spec.Report(report.Terminal{}))
}
