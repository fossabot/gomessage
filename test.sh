#!/bin/bash
set -e

# Run go tests
echo "" > coverage.txt
for package in $(go list ./...); do
    if [ "$(uname -m)" = "x86_64" ]; then
	pushd $(basename $package) > /dev/null
        go test -mod=vendor -tags "$BUILD_TAGS" -race -coverprofile=../profile.out -covermode=atomic $package
	popd > /dev/null
    else
	pushd $(basename $package) > /dev/null
        go test -mod=vendor -tags "$BUILD_TAGS" -coverprofile=../profile.out $package
	popd > /dev/null
    fi

    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
