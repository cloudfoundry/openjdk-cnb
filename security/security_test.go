/*
 * Copyright 2018-2019, Pivotal Software, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of Pivotal Software,
 * Inc.
 */

package security_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/pivotal-cf/p-zulu-cnb/security"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestSecurity(t *testing.T) {
	spec.Run(t, "Java Security Properties", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("contributes java security properties", func() {
			f := test.NewBuildFactory(t)

			s := security.NewSecurity(f.Build)

			g.Expect(s.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer("java-security-properties")
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))

			destination := filepath.Join(layer.Root, "java.security")
			g.Expect(destination).To(gomega.BeARegularFile())
			g.Expect(layer).To(test.HaveAppendLaunchEnvironment("JAVA_OPTS", " -Djava.security.properties=%s", destination))
		})
	}, spec.Report(report.Terminal{}))
}
