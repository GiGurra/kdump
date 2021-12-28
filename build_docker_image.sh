#!/usr/bin/env bash

set -e

PACKAGE_VERSION=$(cat VERSION)

echo "Building kdump version '$PACKAGE_VERSION'"
go build -ldflags="-X 'main.Version=$PACKAGE_VERSION'" .

DOCKER_BUILDKIT=1 docker build . -t gigurra/kdump:$PACKAGE_VERSION
