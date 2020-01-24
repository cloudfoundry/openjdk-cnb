# `openjdk-cnb`
The Cloud Foundry OpenJDK Buildpack is a Cloud Native Buildpack V3 that provides OpenJDK JREs and JDKs to applications.

This buildpack is designed to work in collaboration with other buildpacks which request contributions of JREs and JDKs.

## Behavior
This buildpack will participate if any of the following conditions are met

* Another buildpack requires `openjdk-jdk`
* Another buildpack requires `openjdk-jre`

The buildpack will do the following if a JDK is requested:

* Contributes a JDK to a layer marked `build` and `cache` with all commands on `$PATH`
* Contributes `$JAVA_HOME` configured to the build layer
* Contributes `$JDK_HOME` configure to the build layer

The buildpack will do the following if a JRE is requested:

* Contributes a JRE to a layer with all commands on `$PATH`
* Contributes `$JAVA_HOME` configured to the layer
* Contributes `-XX:ActiveProcessorCount` to the layer
* Contributes `$MALLOC_ARENA_MAX` to the layer
* Disables JVM DNS caching if link-local DNS is available
* If `metadata.build = true`
  * Marks layer as `build` and `cache`
* If `metadata.launch = true`
  * Marks layer as `launch`
* Contributes `jvmkill` to a layer marked `launch`
* Contributes Memory Calculator to a layer marked `launch`

## Configuration
| Environment Variable | Description
| -------------------- | -----------
| `$BP_JAVA_VERSION` | Configure a specific JDK or JRE version.  This value must _exactly_ match a version available in the buildpack so typically it would configured to a wildcard such as `8.*`.
| `$BPL_HEAD_ROOM` | Configure the percentage of headroom the memory calculator will allocated.  Defaults to `0`.
| `$BPL_LOADED_CLASS_COUNT` | Configure the number of classes that will be loaded at runtime.  Defaults to 35% of the number of classes.
| `$BPL_THREAD_COUNT` | Configure the number of user threads at runtime.  Defaults to `250`.


## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: https://www.apache.org/licenses/LICENSE-2.0
