#!/usr/bin/env bash

set -e

export VERSION=$(cat version.go | grep Version | awk '{print $4}' | tr -d '"')
export DOCKER_TAG=gigurra/kdump:$VERSION

echo "Building kdump $VERSION"
go build .

