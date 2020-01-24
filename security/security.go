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

package security

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// Security represents the layer that will hold a java.security properties file at runtime.
type Security struct {
	buildpack buildpack.Buildpack
	layer     layers.Layer
}

// Contribute makes the contribution to launch.
func (s Security) Contribute() error {
	return s.layer.Contribute(marker{s.buildpack.Info, "Java Security Properties"}, func(layer layers.Layer) error {
		if err := os.RemoveAll(layer.Root); err != nil {
			return err
		}

		if err := helper.WriteFile(s.Target(), 0644, ""); err != nil {
			return err
		}

		return layer.AppendLaunchEnv("JAVA_OPTS", " -Djava.security.properties=%s", s.Target())
	}, layers.Launch)
}

// Target returns the target file name representing the java.security file.
func (s Security) Target() string {
	return filepath.Join(s.layer.Root, "java.security")
}

// NewSecurity creates a new Security instance.
func NewSecurity(build build.Build) Security {
	return Security{
		build.Buildpack,
		build.Layers.Layer("java-security-properties"),
	}
}

type marker struct {
	buildpack.Info

	DisplayName string `toml:"display_name"`
}

func (m marker) Identity() (string, string) {
	return "Java Security Properties", m.Version
}
