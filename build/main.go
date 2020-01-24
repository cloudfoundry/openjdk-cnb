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
	"fmt"
	"os"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/openjdk-cnb/dns"
	"github.com/cloudfoundry/openjdk-cnb/jdk"
	"github.com/cloudfoundry/openjdk-cnb/jre"
	"github.com/cloudfoundry/openjdk-cnb/jvmkill"
	"github.com/cloudfoundry/openjdk-cnb/memcalc"
	"github.com/cloudfoundry/openjdk-cnb/provider"
	"github.com/cloudfoundry/openjdk-cnb/security"
)

func main() {
	build, err := build.DefaultBuild()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Build: %s\n", err)
		os.Exit(101)
	}

	if code, err := b(build); err != nil {
		build.Logger.TerminalError(build.Buildpack, err.Error())
		os.Exit(code)
	} else {
		os.Exit(code)
	}
}

func b(build build.Build) (int, error) {
	build.Logger.Title(build.Buildpack)

	if jdk, ok, err := jdk.NewJDK(build); err != nil {
		return build.Failure(102), err
	} else if ok {
		if err := jdk.Contribute(); err != nil {
			return build.Failure(103), err
		}
	}

	if jre, ok, err := jre.NewJRE(build); err != nil {
		return build.Failure(102), err
	} else if ok {
		if err := jre.Contribute(); err != nil {
			return build.Failure(103), err
		}

		s := security.NewSecurity(build)
		if err := s.Contribute(); err != nil {
			return build.Failure(103), err
		}

		if p, ok, err := provider.NewSecurityProviderConfigurer(build, jre.Dependency, s.Target()); err != nil {
			return build.Failure(102), err
		} else if ok {
			if err := p.Contribute(); err != nil {
				return build.Failure(103), err
			}
		}

		if err := dns.NewLinkLocalDNS(build, s.Target()).Contribute(); err != nil {
			return build.Failure(103), err
		}

		if kill, err := jvmkill.NewJVMKill(build); err != nil {
			return build.Failure(102), err
		} else {
			if err := kill.Contribute(); err != nil {
				return build.Failure(103), err
			}
		}

		if err := memcalc.NewClassCounter(build).Contribute(); err != nil {
			return build.Failure(103), err
		}

		if mc, err := memcalc.NewMemoryCalculator(build); err != nil {
			return build.Failure(102), err
		} else {
			if err := mc.Contribute(); err != nil {
				return build.Failure(103), err
			}
		}
	}

	return build.Success()
}
