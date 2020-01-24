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

package jvmkill

import (
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// Dependency is a build plan dependency indicating a requirement for the JVMKill utility.
const Dependency string = "jvmkill"

// JVMKill represents the jvmkill configuration for a JVM application
type JVMKill struct {
	layer layers.DependencyLayer
}

// Contribute makes the contribution to launch
func (j JVMKill) Contribute() error {
	return j.layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.Body("Copying to %s", layer.Root)

		destination := filepath.Join(layer.Root, layer.ArtifactName())

		if err := helper.CopyFile(artifact, destination); err != nil {
			return err
		}

		return layer.AppendSharedEnv("JAVA_OPTS", " -agentpath:%s=printHeapHistogram=1", destination)
	}, layers.Launch)
}

// NewJVMKill creates a new JVMKill instance.
func NewJVMKill(build build.Build) (JVMKill, error) {
	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return JVMKill{}, err
	}

	dep, err := deps.Best(Dependency, "", build.Stack)
	if err != nil {
		return JVMKill{}, err
	}

	return JVMKill{build.Layers.DependencyLayer(dep)}, nil
}
