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

package memcalc_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/cloudfoundry/openjdk-cnb/memcalc"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestClassCounter(t *testing.T) {
	spec.Run(t, "Class Counter", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("contributes class counter", func() {
			f := test.NewBuildFactory(t)
			test.TouchFile(t, f.Build.Buildpack.Root, "bin", "class-counter")

			c := memcalc.NewClassCounter(f.Build)

			g.Expect(c.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("class-counter")
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
			g.Expect(filepath.Join(layer.Root, "bin", "class-counter")).To(gomega.BeARegularFile())
		})
	}, spec.Report(report.Terminal{}))
}
