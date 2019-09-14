version() {
  local PATTERN="([0-9]+)\.([0-9]+)\.([0-9]+)\+(.*)"

  for VERSION in $(cat "../dependency/version"); do
      if [[ ${VERSION} =~ ${PATTERN} ]]; then
        echo "${BASH_REMATCH[1]}.${BASH_REMATCH[2]}.${BASH_REMATCH[3]}"
        return
      else
        >2 echo "version is not semver"
        exit 1
      fi
    done
}
