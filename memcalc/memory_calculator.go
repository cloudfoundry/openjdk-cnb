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

	"github.com/buildpacks/libbuildpack/v2/application"
	"github.com/cloudfoundry/libcfbuildpack/v2/build"
	"github.com/cloudfoundry/libcfbuildpack/v2/helper"
	"github.com/cloudfoundry/libcfbuildpack/v2/layers"
)

// Dependency is a build plan dependency indicating a requirement for the Memory Calculator utility.
const Dependency string = "memory-calculator"

// MemoryCalculator represents the memory calculator configuration for a JVM application.
type MemoryCalculator struct {
	application application.Application
	layer       layers.DependencyLayer
}

// Contribute makes the contribution to launch
func (m MemoryCalculator) Contribute() error {
	return m.layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.LaunchConfiguration("Set $BPL_HEAD_ROOM to configure", "0")
		layer.Logger.LaunchConfiguration("Set $BPL_LOADED_CLASS_COUNT to configure", "35% of classes")
		layer.Logger.LaunchConfiguration("Set $BPL_THREAD_COUNT to configure", "250")

		layer.Logger.Body("Expanding to %s", layer.Root)

		if err := helper.ExtractTarGz(artifact, filepath.Join(layer.Root, "bin"), 0); err != nil {
			return err
		}

		return layer.WriteProfile("memory-calculator", `HEAD_ROOM=${BPL_HEAD_ROOM:=0}

if [[ -z "${BPL_LOADED_CLASS_COUNT+x}" ]]; then
    LOADED_CLASS_COUNT=$(class-counter %s %s)
else
	LOADED_CLASS_COUNT=${BPL_LOADED_CLASS_COUNT}
fi

THREAD_COUNT=${BPL_THREAD_COUNT:=250}

TOTAL_MEMORY=$(cat /sys/fs/cgroup/memory/memory.limit_in_bytes)

if [ ${TOTAL_MEMORY} -eq 9223372036854771712 ]; then
  printf "Container memory limit unset. Configuring JVM for 1G container.\n"
  TOTAL_MEMORY=1073741824
elif [ ${TOTAL_MEMORY} -gt 70368744177664 ]; then
  printf "Container memory limit too large. Configuring JVM for 64T container.\n"
  TOTAL_MEMORY=70368744177664
fi

MEMORY_CONFIGURATION=$(java-buildpack-memory-calculator \
    --head-room "${HEAD_ROOM}" \
    --jvm-options "${JAVA_OPTS}" \
    --loaded-class-count "${LOADED_CLASS_COUNT}" \
    --thread-count "${THREAD_COUNT}" \
    --total-memory "${TOTAL_MEMORY}")

printf "Calculated JVM Memory Configuration: ${MEMORY_CONFIGURATION} (Head Room: ${HEAD_ROOM}%%%%, Loaded Class Count: ${LOADED_CLASS_COUNT}, Thread Count: ${THREAD_COUNT}, Total Memory: ${TOTAL_MEMORY})\n"
export JAVA_OPTS="${JAVA_OPTS} ${MEMORY_CONFIGURATION}"
`, m.application.Root, layer.Dependency.Version.Version.Original())
	}, layers.Launch)
}

// NewMemoryCalculator creates a new MemoryCalculator instance.
func NewMemoryCalculator(build build.Build) (MemoryCalculator, error) {
	deps, err := build.Buildpack.Dependencies()
	if err != nil {
		return MemoryCalculator{}, err
	}

	dep, err := deps.Best(Dependency, "", build.Stack)
	if err != nil {
		return MemoryCalculator{}, err
	}

	return MemoryCalculator{build.Application, build.Layers.DependencyLayer(dep)}, nil
}
