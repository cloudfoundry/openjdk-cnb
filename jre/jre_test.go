/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jre_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/openjdk-buildpack/jre"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestJRE(t *testing.T) {
	spec.Run(t, "JRE", testJRE, spec.Report(report.Terminal{}))
}

func testJRE(t *testing.T, when spec.G, it spec.S) {

	it("returns true if build plan exists", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, jre.Dependency, "stub-openjdk-jre.tar.gz")
		f.AddBuildPlan(t, jre.Dependency, buildplan.Dependency{})

		_, ok, err := jre.NewJRE(f.Build)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("NewJRE = %t, expected true", ok)
		}
	})

	it("returns false if build plan does not exist", func() {
		f := test.NewBuildFactory(t)

		_, ok, err := jre.NewJRE(f.Build)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Errorf("NewJRE = %t, expected false", ok)
		}
	})

	it("contributes JRE to build", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, jre.Dependency, "stub-openjdk-jre.tar.gz")
		f.AddBuildPlan(t, jre.Dependency, buildplan.Dependency{
			Metadata: buildplan.Metadata{jre.BuildContribution: true},
		})

		j, _, err := jre.NewJRE(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := j.Contribute(); err != nil {
			t.Fatal(err)
		}

		layer := f.Build.Layers.Layer("openjdk-jre")
		test.BeLayerLike(t, layer, true, true, false)
		test.BeFileLike(t, filepath.Join(layer.Root, "fixture-marker"), 0644, "")
		test.BeOverrideSharedEnvLike(t, layer, "JAVA_HOME", layer.Root)
	})

	it("contributes JRE to launch", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, jre.Dependency, "stub-openjdk-jre.tar.gz")
		f.AddBuildPlan(t, jre.Dependency, buildplan.Dependency{
			Metadata: buildplan.Metadata{jre.LaunchContribution: true},
		})

		j, _, err := jre.NewJRE(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := j.Contribute(); err != nil {
			t.Fatal(err)
		}

		layer := f.Build.Layers.Layer("openjdk-jre")
		test.BeLayerLike(t, layer, false, true, true)
		test.BeFileLike(t, filepath.Join(layer.Root, "fixture-marker"), 0644, "")
		test.BeOverrideSharedEnvLike(t, layer, "JAVA_HOME", layer.Root)
	})
}
