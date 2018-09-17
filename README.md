# `openjdk-buildpack`
The Cloud Foundry OpenJDK Buildpack is a Cloud Native Buildpack V3 that provides OpenJDK JREs and JDKs to applications.

## Detection
Detection always passes and contributes nothing to the build plan.

## Build

```toml
[openjdk-jdk]
```

* Contributes a JDK to a cache layer named the same as the key.
* Contributes `JAVA_HOME` configured to a cache layer named the same as the key.
* Contributes `JDK_HOME` configured to a cache layer named the same as the key.

```toml
[openjdk-jre]

  [openjdk-jre.metadata]
  build  = true
  launch = true
```

* Contributes a JRE to a cache layer named the same as the key if `build` metadata key is set.
* Contributes `JAVA_HOME` configured to a cache layer named the same as the key if `build` metadata key is set.
* Contributes a JRE to a launch layer if `launch` metadata key is set.
* Contributes `.profile.d/java-home` if `launch` metadata key is set.

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

