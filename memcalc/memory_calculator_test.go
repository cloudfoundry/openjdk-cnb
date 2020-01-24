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

func TestMemoryCalculator(t *testing.T) {
	spec.Run(t, "Memory Calculator", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("contributes memory calculator", func() {
			f := test.NewBuildFactory(t)
			f.AddDependency(memcalc.Dependency, filepath.Join("testdata", "stub-memory-calculator.tgz"))

			m, err := memcalc.NewMemoryCalculator(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(m.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("memory-calculator")
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
			g.Expect(filepath.Join(layer.Root, "bin", "java-buildpack-memory-calculator")).To(gomega.BeARegularFile())
			g.Expect(layer).To(test.HaveProfile("memory-calculator", `HEAD_ROOM=${BPL_HEAD_ROOM:=0}

if [[ -z "${BPL_LOADED_CLASS_COUNT+x}" ]]; then
    LOADED_CLASS_COUNT=$(class-counter %s 1.0)
else
	LOADED_CLASS_COUNT=${BPL_LOADED_CLASS_COUNT}
fi

THREAD_COUNT=${BPL_THREAD_COUNT:=250}

TOTAL_MEMORY=$(cat /sys/fs/cgroup/memory/memory.limit_in_bytes)
TOTAL_MEMORY=$((TOTAL_MEMORY < 70368744177664 ? TOTAL_MEMORY : 70368744177664))

MEMORY_CONFIGURATION=$(java-buildpack-memory-calculator \
    --head-room "${HEAD_ROOM}" \
    --jvm-options "${JAVA_OPTS}" \
    --loaded-class-count "${LOADED_CLASS_COUNT}" \
    --thread-count "${THREAD_COUNT}" \
    --total-memory "${TOTAL_MEMORY}")

printf "Calculated JVM Memory Configuration: ${MEMORY_CONFIGURATION} (Head Room: ${HEAD_ROOM}%%%%, Loaded Class Count: ${LOADED_CLASS_COUNT}, Thread Count: ${THREAD_COUNT}, Total Memory: ${TOTAL_MEMORY})\n"
export JAVA_OPTS="${JAVA_OPTS} ${MEMORY_CONFIGURATION}"
`, f.Build.Application.Root))
		})
	}, spec.Report(report.Terminal{}))
}
