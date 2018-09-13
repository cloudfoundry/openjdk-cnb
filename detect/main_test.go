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

package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", testDetect, spec.Report(report.Terminal{}))
}

func testDetect(t *testing.T, when spec.G, it spec.S) {

	it("always passes", func() {
		root := libjavabuildpack.ScratchDir(t, "detect")
		defer libjavabuildpack.ReplaceWorkingDirectory(t, root)()

		defer libjavabuildpack.ReplaceEnv(t, "PACK_STACK_ID", "test-stack")()

		c, d := libjavabuildpack.ReplaceConsole(t)
		defer d()
		c.In(t, "")

		err := libjavabuildpack.WriteToFile(strings.NewReader(""), filepath.Join(root, "buildpack.toml"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		defer libjavabuildpack.ReplaceArgs(t, filepath.Join(root, "bin", "test"))()

		actual, d := libjavabuildpack.CaptureExitStatus(t)
		defer d()

		main()

		if *actual != 0 {
			t.Errorf("os.Exit = %d, expected 0", *actual)
		}

		actualStdout := c.Out(t)
		if actualStdout != "" {
			t.Errorf("stdout = %s, expected empty", actualStdout)
		}
	})

}
