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

package provider

import (
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/v2/build"
	"github.com/cloudfoundry/libcfbuildpack/v2/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/v2/helper"
	"github.com/cloudfoundry/libcfbuildpack/v2/layers"
)

// Security Provider represents the SecurityProviderConfigurer for a JVM application.
type SecurityProviderConfigurer struct {
	buildpack buildpack.Buildpack
	jre       buildpack.Dependency
	layer     layers.HelperLayer
	security  string
}

// Contribute makes the contribution to launch.
func (s SecurityProviderConfigurer) Contribute() error {
	return s.layer.Contribute(func(artifact string, layer layers.HelperLayer) error {
		if err := helper.CopyFile(artifact, filepath.Join(layer.Root, "bin", "security-provider-configurer")); err != nil {
			return err
		}

		if err := s.classpath(layer); err != nil {
			return err
		}

		return layer.WriteProfile("security-provider-configurer",
			`security-provider-configurer $JAVA_HOME %s %s $SECURITY_PROVIDERS`, s.jre.Version.Original(), s.security)
	}, layers.Launch)
}

func (s SecurityProviderConfigurer) classpath(layer layers.HelperLayer) error {
	if helper.AfterJava8(s.jre.Version.Version) {
		return layer.WriteProfile("security-provider-classpath", "export CLASSPATH=$CLASSPATH:$SECURITY_PROVIDERS_CLASSPATH")
	}

	return layer.WriteProfile("security-provider-classpath", `EXT_DIRS="$JAVA_HOME/lib/ext"

for I in ${SECURITY_PROVIDERS_CLASSPATH//:/$'\n'}; do
  EXT_DIRS="$EXT_DIRS:$(dirname $I)"
done

JAVA_OPTS="$JAVA_OPTS -Djava.ext.dirs=$EXT_DIRS"`)
}

// NewSecurityProviderConfigurer creates a new SecurityProviderConfigurer instance.
func NewSecurityProviderConfigurer(build build.Build, jre buildpack.Dependency, security string) (SecurityProviderConfigurer, bool, error) {
	return SecurityProviderConfigurer{
		build.Buildpack,
		jre,
		build.Layers.HelperLayer("security-provider-configurer", "Security Provider Configurer"),
		security,
	}, true, nil
}
