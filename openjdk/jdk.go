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

package openjdk

import (
	"fmt"
	"reflect"

	"github.com/buildpacks/libcnb"
	"github.com/paketoio/libpak"
	"github.com/paketoio/libpak/bard"
	"github.com/paketoio/libpak/crush"
)

type JDK struct {
	Crush           crush.Crush
	Dependency      libpak.BuildpackDependency
	DependencyCache libpak.DependencyCache
	Logger          bard.Logger
}

func (j JDK) Name() string {
	return "jdk"
}

func (j JDK) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	if reflect.DeepEqual(j.Dependency.Metadata(), layer.Metadata) {
		return layer, nil
	}

	artifact, err := j.DependencyCache.Artifact(j.Dependency)
	if err != nil {
		return libcnb.Layer{}, fmt.Errorf("unable to get artifact: %w", err)
	}
	defer artifact.Close()

	j.Logger.Body("Expanding to %s", layer.Path)
	err = j.Crush.ExtractTarGz(artifact, layer.Path, 1)
	if err != nil {
		return libcnb.Layer{}, fmt.Errorf("unable to expand JDK: %w", err)
	}

	layer.LaunchEnvironment.Override("JAVA_HOME", layer.Path)
	layer.LaunchEnvironment.Override("JDK_HOME", layer.Path)

	layer.Build = true
	layer.Launch = true
	layer.Metadata = j.Dependency.Metadata()

	return layer, nil
}
