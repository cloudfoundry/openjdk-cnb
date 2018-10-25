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

package openjdk_buildpack

import (
	"fmt"

	"github.com/cloudfoundry/libjavabuildpack"
)

const (
	// BuildContribution is a build plan dependency key indicating a requirement for the dependency at build time.
	BuildContribution string = "build"

	// JREDependency is a build plan dependency indicating a requirement for a JRE.
	JREDependency string = "openjdk-jre"

	// LaunchContribution is a build plan dependency yet indicate a requirement for the dependency at launch time.
	LaunchContribution string = "launch"
)

// JRE represents a JRE contribution by the buildpack.
type JRE struct {
	buildContribution  bool
	buildLayer         libjavabuildpack.DependencyCacheLayer
	launchContribution bool
	launchLayer        libjavabuildpack.DependencyLaunchLayer
}

// Contribute contributes an expanded JRE to a cache layer.
func (j JRE) Contribute() error {
	if j.buildContribution {
		return j.buildLayer.Contribute(func(artifact string, layer libjavabuildpack.DependencyCacheLayer) error {
			layer.Logger.SubsequentLine("Expanding to %s", layer.Root)
			if err := libjavabuildpack.ExtractTarGz(artifact, layer.Root, 0) ; err != nil {
				return err
			}

			layer.OverrideEnv("JAVA_HOME", layer.Root)

			return nil
		})
	}

	if j.launchContribution {
		return j.launchLayer.Contribute(func(artifact string, layer libjavabuildpack.DependencyLaunchLayer) error {
			layer.Logger.SubsequentLine("Expanding to %s", layer.Root)
			if err := libjavabuildpack.ExtractTarGz(artifact, layer.Root, 0) ; err != nil {
				return err
			}

			layer.WriteProfile("java-home", "export JAVA_HOME=%s", layer.Root)

			return nil
		})
	}

	return nil
}

// String makes JRE satisfy the Stringer interface.
func (j JRE) String() string {
	return fmt.Sprintf("JRE{ buildContribution: %t, buildLayer: %s, launchContribution: %t, launchLayer: %s }",
		j.buildContribution, j.buildLayer, j.launchContribution, j.launchLayer)
}

// NewJRE creates a new JRE instance. OK is true if build plan contains "openjdk-jre" dependency, otherwise false.
func NewJRE(build libjavabuildpack.Build) (JRE, bool, error) {
	bp, ok := build.BuildPlan[JREDependency]
	if !ok {
		return JRE{}, false, nil
	}

	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return JRE{}, false, err
	}

	dep, err := deps.Best(JREDependency, bp.Version, build.Stack)
	if err != nil {
		return JRE{}, false, err
	}

	jre := JRE{}

	if _, ok := bp.Metadata[BuildContribution]; ok {
		jre.buildContribution = true
		jre.buildLayer = build.Cache.DependencyLayer(dep)
	}

	if _, ok := bp.Metadata[LaunchContribution]; ok {
		jre.launchContribution = true
		jre.launchLayer = build.Launch.DependencyLayer(dep)
	}

	return jre, true, nil
}
