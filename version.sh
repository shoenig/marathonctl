#!/bin/bash

set -euo pipefail

case ${1} in
    version)
        hash=$(git rev-parse HEAD)
        sed -i s/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx/${hash}/ version.go
        ;;
    semver)
        latest=$(git tag --sort=version:refname | tail -n 1)
        sed -i s/x.x.x/${latest}/ version.go
        ;;
    *)
        echo "unknown flag: ${1}"
        exit 1
        ;;
esac