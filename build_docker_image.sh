#!/usr/bin/env bash

set -e

echo "Building kdump"
go build .

PACKAGE_VERSION=$(./kdump -v | awk 'NF>1{print $NF}')

echo "Built kdump version: $PACKAGE_VERSION"

DOCKER_BUILDKIT=1 docker build . -t gigurra/kdump:$PACKAGE_VERSION
