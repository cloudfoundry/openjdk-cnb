#!/usr/bin/env bash

set -euo pipefail

if [[ -d $PWD/go-module-cache && ! -d ${GOPATH}/pkg/mod ]]; then
  mkdir -p ${GOPATH}/pkg
  ln -s $PWD/go-module-cache ${GOPATH}/pkg/mod
fi

version() {
  local PATTERN="([0-9]+)\.([0-9]+)\.([0-9]+)-?(.*)"

  for VERSION in $(cat $1); do
      if [[ ${VERSION} =~ ${PATTERN} ]]; then
        if [[ ${BASH_REMATCH[2]} == "8" ]]; then
          echo "${BASH_REMATCH[1]}.${BASH_REMATCH[2]}.${BASH_REMATCH[4]}"
        else
          echo "${BASH_REMATCH[1]}.${BASH_REMATCH[2]}.${BASH_REMATCH[3]}"
        fi

        return
      else
        >2 echo "version is not semver"
        exit 1
      fi
    done
}

cd "$(dirname "${BASH_SOURCE[0]}")/.."

go build -ldflags='-s -w' -o bin/dependency github.com/cloudfoundry/libcfbuildpack/dependency

bin/dependency openjdk-jre "1\.8\.[\d]+"  $(version ../openjdk-jre-8/version)  $(cat ../openjdk-jre-8/uri)  $(cat ../openjdk-jre-8/sha256)
bin/dependency openjdk-jdk "1\.8\.[\d]+"  $(version ../openjdk-jdk-8/version)  $(cat ../openjdk-jdk-8/uri)  $(cat ../openjdk-jdk-8/sha256)

bin/dependency openjdk-jre "11\.0\.[\d]+" $(version ../openjdk-jre-11/version) $(cat ../openjdk-jre-11/uri) $(cat ../openjdk-jre-11/sha256)
bin/dependency openjdk-jdk "11\.0\.[\d]+" $(version ../openjdk-jdk-11/version) $(cat ../openjdk-jdk-11/uri) $(cat ../openjdk-jdk-11/sha256)

bin/dependency openjdk-jre "12\.0\.[\d]+" $(version ../openjdk-jre-12/version) $(cat ../openjdk-jre-12/uri) $(cat ../openjdk-jre-12/sha256)
bin/dependency openjdk-jdk "12\.0\.[\d]+" $(version ../openjdk-jdk-12/version) $(cat ../openjdk-jdk-12/uri) $(cat ../openjdk-jdk-12/sha256)

cat buildpack.toml
