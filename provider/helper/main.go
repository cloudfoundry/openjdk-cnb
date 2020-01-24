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
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/Masterminds/semver"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/magiconair/properties"
)

func main() {
	root := os.Args[1]

	version, err := semver.NewVersion(os.Args[2])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	security := os.Args[3]

	code, err := p(security, root, version, os.Args[4:])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(code)
	}

	os.Exit(code)
}

func p(out string, root string, version *semver.Version, additional []string) (int, error) {
	if len(additional) == 0 {
		return 0, nil
	}

	fmt.Println("Adding Security Providers to JVM")

	in, err := security(root, version)
	if err != nil {
		return 1, err
	}

	providers, err := readProviders(in)
	if err != nil {
		return 1, err
	}

	providers, err = addProviders(providers, additional)
	if err != nil {
		return 1, err
	}

	if err := writeProviders(providers, out); err != nil {
		return 1, err
	}

	return 0, nil
}

func addProviders(base []string, additional []string) ([]string, error) {
	p := base

	r := regexp.MustCompile(`(?:([\d]+)\|)?([\w.]+)`)
	for _, a := range additional {
		matches := r.FindStringSubmatch(a)

		if matches[1] == "" {
			p = append(p, matches[2])
			continue
		}

		i, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}

		p = append(p, "")
		copy(p[i:], p[i-1:])
		p[i-1] = matches[2]
	}

	return p, nil
}

func readProviders(in string) ([]string, error) {
	if exists, err := helper.FileExists(in); err != nil {
		return nil, err
	} else if !exists {
		return []string{}, nil
	}

	p, err := properties.LoadFile(in, properties.UTF8)
	if err != nil {
		return nil, err
	}
	p = p.FilterStripPrefix("security.provider.")

	keys := p.Keys()
	sort.Slice(keys, func(i, j int) bool {
		a, err := strconv.Atoi(keys[i])
		if err != nil {
			return false
		}

		b, err := strconv.Atoi(keys[j])
		if err != nil {
			return false
		}

		return b > a
	})

	var providers []string

	for _, q := range keys {
		providers = append(providers, p.MustGet(q))
	}

	return providers, nil
}

func security(root string, version *semver.Version) (string, error) {
	if helper.BeforeJava9(version) {
		return filepath.Join(root, "lib", "security", "java.security"), nil
	}

	return filepath.Join(root, "conf", "security", "java.security"), nil
}

func writeProviders(providers []string, out string) error {
	if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("\n"); err != nil {
		return err
	}

	for i, p := range providers {
		if _, err := f.WriteString(fmt.Sprintf("security.provider.%d=%s\n", i+1, p)); err != nil {
			return err
		}
	}

	return nil
}
