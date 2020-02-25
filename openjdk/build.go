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

	"github.com/buildpacks/libcnb"
	"github.com/paketoio/libpak"
	"github.com/paketoio/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	md, err := libpak.NewBuildpackMetadata(context.Buildpack.Metadata)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to unmarshal buildpack metadata: %w", err)
	}

	dr := libpak.DependencyResolver{Dependencies: md.Dependencies}

	b.Logger.Title(context.Buildpack)
	result := libcnb.BuildResult{}

	_, ok, err := pr.Resolve("jdk")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve jdk plan entry: %w", err)
	}

	if ok {
		// TODO: Determine version
		dep, err := dr.Resolve(libpak.DependencyConstraint{ID: "jdk", StackID: context.StackID})
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to resolve jdk dependency: %w", err)
		}
		result.Layers = append(result.Layers, JDK{Dependency: dep, Logger: b.Logger})
	}

	jre, ok, err := pr.Resolve("jre")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve jre plan entry: %w", err)
	}

	if ok {
		// TODO: Determine version
		dep, err := dr.Resolve(libpak.DependencyConstraint{ID: "jre", StackID: context.StackID})
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to resolve jre dependency: %w", err)
		}
		// TODO: layer flags error handling
		result.Layers = append(result.Layers, JRE{Dependency: dep, Logger: b.Logger, Build: jre.Metadata["build"].(bool), Launch: jre.Metadata["launch"].(bool)})
	}

	// if jdk, jdk layer
	// if jre, jre layer
	//   jre metadata describes launch and build

	// TODO: BuildPlan

	// dep := libpak.BuildpackDependency{}

	// return libcnb.BuildResult{
	// 	Layers: []libcnb.LayerContributor{
	// 		JDK{Helper: DependencyHelper{dep}, Logger: b.Logger},
	// 	},
	// }, nil

	// if jdk, ok, err := pr.Resolve("jdk"); err != nil {
	// 	return result, fmt.Errorf("unable to resolve buildpack plan entry jdk: %w", err)
	// } else if ok {
	// 	v := VersionResolver{
	// 		Entry:           jdk,
	// 		DefaultVersions: context.Buildpack.Metadata.DefaultVersions,
	// 	}
	//
	// 	pr := libpak.DependencyResolver{Dependencies: context.Buildpack.Metadata.Dependencies}
	//
	// 	constraint := libpak.DependencyConstraint{
	// 		ID:      "jdk",
	// 		Version: v.Resolve("jdk", "BP_JAVA_VERSION"),
	// 		StackID: context.StackID,
	// 	}
	//
	// 	dep, err := pr.Resolve(constraint)
	// 	if err != nil {
	// 		return result, fmt.Errorf("unable to resolve jdk dependency in %s: %w", context.Buildpack.Metadata.Dependencies, err)
	// 	}
	//
	// }
	//
	// if _, ok, err := pr.Resolve("jre"); err != nil {
	// 	return result, fmt.Errorf("unable to resolve buildpack plan entry jre: %w", err)
	// } else if ok {
	//
	// }

	return result, nil
}
