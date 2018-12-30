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

package jdk

import (
	"fmt"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/openjdk-buildpack/internal"
)

// Dependency is a build plan dependency indicating a requirement for a JDK.
const Dependency = "openjdk-jdk"

// JDK represents a JDK contribution by the buildpack.
type JDK struct {
	layer layers.DependencyLayer
}

// Contribute contributes an expanded JDK to a cache layer.
func (j JDK) Contribute() error {
	return j.layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.SubsequentLine("Expanding to %s", layer.Root)

		if err := helper.ExtractTarGz(artifact, layer.Root, 0); err != nil {
			return err
		}

		if err := layer.OverrideBuildEnv("JAVA_HOME", layer.Root); err != nil {
			return err
		}

		if err := layer.OverrideBuildEnv("JDK_HOME", layer.Root); err != nil {
			return err
		}

		return nil
	}, layers.Build, layers.Cache)
}

// String makes JDK satisfy the Stringer interface.
func (j JDK) String() string {
	return fmt.Sprintf("JDK{ layer: %s }", j.layer)
}

// NewJDK creates a new JDK instance. OK is true if build plan contains "openjdk-jdk" dependency, otherwise false.
func NewJDK(build build.Build) (JDK, bool, error) {
	bp, ok := build.BuildPlan[Dependency]
	if !ok {
		return JDK{}, false, nil
	}

	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return JDK{}, false, err
	}

	version := internal.Version(build.Buildpack, bp)

	dep, err := deps.Best(Dependency, version, build.Stack)
	if err != nil {
		return JDK{}, false, err
	}

	return JDK{build.Layers.DependencyLayer(dep)}, true, nil
}
