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

package openjdk_buildpack_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/cloudfoundry/libjavabuildpack/test"
	"github.com/cloudfoundry/openjdk-buildpack"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestJRE(t *testing.T) {
	spec.Run(t, "JRE", testJRE, spec.Report(report.Terminal{}))
}

func testJRE(t *testing.T, when spec.G, it spec.S) {

	it("returns true if build plan exists", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, openjdk_buildpack.JREDependency, "stub-openjdk-jre.tar.gz")
		f.AddBuildPlan(t, openjdk_buildpack.JREDependency, libbuildpack.BuildPlanDependency{})

		_, ok, err := openjdk_buildpack.NewJRE(f.Build)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("NewJRE = %t, expected true", ok)
		}
	})

	it("returns false if build plan does not exist", func() {
		f := test.NewBuildFactory(t)

		_, ok, err := openjdk_buildpack.NewJRE(f.Build)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Errorf("NewJRE = %t, expected false", ok)
		}
	})

	it("contributes JRE to cache", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, openjdk_buildpack.JREDependency, "stub-openjdk-jre.tar.gz")
		f.AddBuildPlan(t, openjdk_buildpack.JREDependency, libbuildpack.BuildPlanDependency{
			Metadata: libbuildpack.BuildPlanDependencyMetadata{openjdk_buildpack.BuildContribution: true},
		})

		j, _, err := openjdk_buildpack.NewJRE(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := j.Contribute(); err != nil {
			t.Fatal(err)
		}

		layerRoot := filepath.Join(f.Build.Cache.Root, "openjdk-jre")
		libjavabuildpack.BeFileLike(t, filepath.Join(layerRoot, "fixture-marker"), 0644, "")
		libjavabuildpack.BeFileLike(t, filepath.Join(layerRoot, "env", "JAVA_HOME.override"), 0644, layerRoot)
	})

	it("contributes JRE to launch", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, openjdk_buildpack.JREDependency, "stub-openjdk-jre.tar.gz")
		f.AddBuildPlan(t, openjdk_buildpack.JREDependency, libbuildpack.BuildPlanDependency{
			Metadata: libbuildpack.BuildPlanDependencyMetadata{openjdk_buildpack.LaunchContribution: true},
		})

		j, _, err := openjdk_buildpack.NewJRE(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := j.Contribute(); err != nil {
			t.Fatal(err)
		}

		layerRoot := filepath.Join(f.Build.Launch.Root, "openjdk-jre")
		libjavabuildpack.BeFileLike(t, filepath.Join(layerRoot, "fixture-marker"), 0644, "")
		libjavabuildpack.BeFileLike(t, filepath.Join(layerRoot, "profile.d", "java-home"), 0644,
			fmt.Sprintf("export JAVA_HOME=%s", layerRoot))
	})
}
