#!/usr/bin/env bash

set -e

export VERSION=$(cat version.go | grep Version | awk '{print $4}' | tr -d '"')

echo "Building kdump $VERSION"
go build .

