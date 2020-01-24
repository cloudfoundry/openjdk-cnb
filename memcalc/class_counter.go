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

package memcalc

import (
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/v2/build"
	"github.com/cloudfoundry/libcfbuildpack/v2/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/v2/helper"
	"github.com/cloudfoundry/libcfbuildpack/v2/layers"
)

// ClassCounter represents the class-counter helper application.
type ClassCounter struct {
	buildpack buildpack.Buildpack
	layer     layers.HelperLayer
}

// Contributes makes the contribution to launch
func (c ClassCounter) Contribute() error {
	return c.layer.Contribute(func(artifact string, layer layers.HelperLayer) error {
		return helper.CopyFile(artifact, filepath.Join(layer.Root, "bin", "class-counter"))
	}, layers.Launch)
}

// NewClassCounter creates a new ClassCounter instance.
func NewClassCounter(build build.Build) ClassCounter {
	return ClassCounter{build.Buildpack, build.Layers.HelperLayer("class-counter", "Class Counter")}
}
