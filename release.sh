#!/bin/bash

set -euo pipefail

previous=$(git tag --sort=version:refname | tail -n 1)

read -p "what is the new release version (previous: ${previous})? " ; new=${REPLY}

# tag this release in git
echo "generating git tag for new version ${new} ..."
git tag "${new}"
git push origin "${new}"

# generate the version.go file
echo "running go generate ..."
go generate

# use gox to build all the cross platform versions
echo "building all versions ..."
gox


#!/bin/bash

# create checksums for each binary
echo "md5 checksums for release ${new}" > CHECKSUM.txt
echo "" >> CHECKSUM.txt
for bin in marathonctl_*; do
        sum=$(md5sum ${bin})
	echo "${sum}" >> CHECKSUM.txt
done

# create a .tar.gz file with all binaries
echo "creating .tar.gz file ..."
targz="marathonctl-${new}.tar.gz"
tar -czvf ${targz} CHECKSUM.txt ./marathonctl_*

# clean up the mess
echo "cleaning up binaries ..."
rm -rf ./marathonctl_*
echo "removing checksum file ..."
rm -rf ./CHECKSUM.txt
echo "resetting version.go ..."
git cout version.go
