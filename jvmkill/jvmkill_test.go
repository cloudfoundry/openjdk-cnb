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

package jvmkill_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/openjdk-cnb/jvmkill"
	"github.com/onsi/gomega"

	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestJVMKill(t *testing.T) {
	spec.Run(t, "JVMKill", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("contributes JVMKill", func() {
			f := test.NewBuildFactory(t)
			f.AddDependency(jvmkill.Dependency, filepath.Join("testdata", "stub-jvmkill.so"))

			j, err := jvmkill.NewJVMKill(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(j.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("jvmkill")
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
			g.Expect(filepath.Join(layer.Root, "stub-jvmkill.so")).To(gomega.BeARegularFile())
			g.Expect(layer).To(test.HaveAppendSharedEnvironment("JAVA_OPTS", " -agentpath:%s/stub-jvmkill.so=printHeapHistogram=1", layer.Root))
		})
	}, spec.Report(report.Terminal{}))
}
