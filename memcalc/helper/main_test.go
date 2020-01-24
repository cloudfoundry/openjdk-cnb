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

package main

import (
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestCount(t *testing.T) {
	spec.Run(t, "Count", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var root string

		it.Before(func() {
			root = test.ScratchDir(t, "count")
		})

		it("counts files in application", func() {
			version, err := semver.NewVersion("8.0.232")
			g.Expect(err).NotTo(gomega.HaveOccurred())

			test.TouchFile(t, root, "alpha.class")
			test.TouchFile(t, root, "bravo", "charlie.class")

			g.Expect(c(root, version)).To(gomega.Equal(27869))
		})

		it("counts files in archives", func() {
			version, err := semver.NewVersion("8.0.232")
			g.Expect(err).NotTo(gomega.HaveOccurred())

			test.CopyFile(t, filepath.Join("testdata", "stub-dependency.jar"), filepath.Join(root, "stub-dependency.jar"))

			g.Expect(c(root, version)).To(gomega.Equal(27869))
		})

		it("counts Java 8 JRE", func() {
			version, err := semver.NewVersion("8.0.232")
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(c(root, version)).To(gomega.Equal(27867))
		})

		it("counts Java 9 JRE", func() {
			version, err := semver.NewVersion("9.0.4")
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(c(root, version)).To(gomega.Equal(25565))
		})

		it("counts Java 10 JRE", func() {
			version, err := semver.NewVersion("10.0.2")
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(c(root, version)).To(gomega.Equal(28191))
		})

		it("counts Java 11 JRE", func() {
			version, err := semver.NewVersion("11.0.1")
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(c(root, version)).To(gomega.Equal(24219))
		})

		it("counts Java 12 JRE", func() {
			version, err := semver.NewVersion("12.0.1")
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(c(root, version)).To(gomega.Equal(24219))
		})

		it("counts Java 13 JRE", func() {
			version, err := semver.NewVersion("12.0.1")
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(c(root, version)).To(gomega.Equal(24219))
		})

		it("counts unknown JRE", func() {
			version, err := semver.NewVersion("14.0.1")
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(c(root, version)).To(gomega.Equal(24219))
		})
	}, spec.Report(report.Terminal{}))
}
