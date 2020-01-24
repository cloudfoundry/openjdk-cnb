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
	"archive/zip"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
)

const loadFactor = 0.35

func main() {
	root := os.Args[1]

	version, err := semver.NewVersion(os.Args[2])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	count, code, err := c(root, version)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(code)
	}

	fmt.Println(math.Ceil(float64(count) * loadFactor))
	os.Exit(code)
}

func c(root string, version *semver.Version) (count int, code int, err error) {
	j, err := jre(version)
	if err != nil {
		return 0, 1, err
	}

	d, err := directory(root)
	if err != nil {
		return 0, 1, err
	}

	return j + d, 0, nil
}

func archive(file string) (int, error) {
	count := 0

	z, err := zip.OpenReader(file)
	if err != nil {
		return 0, err
	}

	for _, f := range z.File {
		if shouldCount(f.Name) {
			count += 1
		}
	}

	return count, nil
}

func directory(root string) (int, error) {
	count := 0

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".jar") {
			a, err := archive(path)
			if err != nil {
				return err
			}

			count += a
			return nil
		}

		if shouldCount(path) {
			count += 1
			return nil
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return count, nil
}

func isVersion(constraint string, version *semver.Version) (bool, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, err
	}

	v, _ := c.Validate(version)
	return v, nil
}

func jre(version *semver.Version) (int, error) {
	if ok, err := isVersion("^8", version); err != nil {
		return 0, err
	} else if ok {
		return 27867, nil
	}

	if ok, err := isVersion("^9", version); err != nil {
		return 0, err
	} else if ok {
		return 25565, nil
	}

	if ok, err := isVersion("^10", version); err != nil {
		return 0, err
	} else if ok {
		return 28191, nil
	}

	if ok, err := isVersion("^11", version); err != nil {
		return 0, err
	} else if ok {
		return 24219, nil
	}

	if ok, err := isVersion("^12", version); err != nil {
		return 0, err
	} else if ok {
		return 24219, nil
	}

	if ok, err := isVersion("^13", version); err != nil {
		return 0, err
	} else if ok {
		return 24219, nil
	}

	return 24219, nil
}

func shouldCount(file string) bool {
	return strings.HasSuffix(file, ".class") ||
		strings.HasSuffix(file, ".groovy") ||
		strings.HasSuffix(file, ".kt")
}
