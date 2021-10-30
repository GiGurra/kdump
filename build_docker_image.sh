#!/usr/bin/env bash

set -e

PACKAGE_VERSION=$(cat package.json \
  | grep version \
  | head -1 \
  | awk -F: '{ print $2 }' \
  | sed 's/[",]//g' \
  | tr -d '[[:space:]]')

echo $PACKAGE_VERSION

DOCKER_BUILDKIT=1 docker build . -t gigurra/kdump:$PACKAGE_VERSION

