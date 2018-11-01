# `openjdk-buildpack`
The Cloud Foundry OpenJDK Buildpack is a Cloud Native Buildpack V3 that provides OpenJDK JREs and JDKs to applications.

This buildpack is designed to work in collaboration with other buildpacks which request contributions of JREs and JDKs.

## Detection
The detection phase always passes and contributes nothing to the build plan, depending on other buildpacks to request contributions.

## Build
If the build plan contains

* `openjdk-jdk`
  * Contributes a JDK to a cache layer with all commands on `$PATH`
  * Contributes `$JAVA_HOME` configured to the cache layer
  * Contributes `$JDK_HOME` configure to the cache layer

* `openjdk-jre`
  * `metadata.build = true`
    * Contributes a JRE to a cache layer with all comands on `$PATH`
    * Contributes `$JAVA_HOME` configured to the cache layer
  * `metadata.launch = true`
    * Contributes a JRE to a launch layer will all commands on `$PATH`
    * Contributes `$JAVA_HOME` configured to the launch layer

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
