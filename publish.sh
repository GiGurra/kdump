#!/usr/bin/env bash

set -e

source build.sh

echo "Publishing kdump $VERSION"

git tag "$VERSION"
git push origin "$VERSION"

docker push "$DOCKER_TAG"
