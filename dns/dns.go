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

package dns

import (
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// LinkLocalDNS represents the Link-Local DNS configurer for a JVM application.
type LinkLocalDNS struct {
	buildpack buildpack.Buildpack
	layer     layers.HelperLayer
	security  string
}

// Contribute makes the contribution to launch.
func (l LinkLocalDNS) Contribute() error {
	return l.layer.Contribute(func(artifact string, layer layers.HelperLayer) error {
		if err := helper.CopyFile(artifact, filepath.Join(layer.Root, "bin", "link-local-dns")); err != nil {
			return err
		}

		return layer.WriteProfile("link-local-dns", `link-local-dns %s`, l.security)
	}, layers.Launch)
}

// NewLinkLocalDNS creates a new LinkLocalDNS instance.
func NewLinkLocalDNS(build build.Build, security string) LinkLocalDNS {
	return LinkLocalDNS{
		build.Buildpack,
		build.Layers.HelperLayer("link-local-dns", "Link-Local DNS"),
		security,
	}
}
