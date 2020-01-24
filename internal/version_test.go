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

package internal_test

import (
	"testing"

	bp "github.com/buildpacks/libbuildpack/v2/buildpack"
	"github.com/buildpacks/libbuildpack/v2/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/v2/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/v2/logger"
	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/cloudfoundry/openjdk-cnb/internal"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestVersion(t *testing.T) {
	spec.Run(t, "Version", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("uses $BP_JAVA_VERSION if set", func() {
			defer test.ReplaceEnv(t, "BP_JAVA_VERSION", "test-version")()
			buildpack := buildpack.NewBuildpack(bp.Buildpack{}, logger.Logger{})
			dependency := buildpackplan.Plan{}

			g.Expect(internal.Version("test-id", dependency, buildpack)).To(gomega.Equal("test-version"))
		})

		it("uses build plan version if set", func() {
			buildpack := buildpack.NewBuildpack(bp.Buildpack{}, logger.Logger{})
			dependency := buildpackplan.Plan{Version: "test-version"}

			g.Expect(internal.Version("test-id", dependency, buildpack)).To(gomega.Equal("test-version"))
		})

		it("uses buildpack default version if set", func() {
			buildpack := buildpack.NewBuildpack(bp.Buildpack{Metadata: buildpack.Metadata{"default-versions": map[string]interface{}{"test-id": "test-version"}}}, logger.Logger{})
			dependency := buildpackplan.Plan{}

			g.Expect(internal.Version("test-id", dependency, buildpack)).To(gomega.Equal("test-version"))
		})

		it("return error if none set", func() {
			buildpack := buildpack.NewBuildpack(bp.Buildpack{Metadata: buildpack.Metadata{"default-versions": map[string]interface{}{"test-id-2": "test-version"}}}, logger.Logger{})
			dependency := buildpackplan.Plan{}

			_, err := internal.Version("test-id", dependency, buildpack)
			g.Expect(err).To(gomega.MatchError("test-id does not map to a string in default-versions map"))
		})
	}, spec.Report(report.Terminal{}))
}
