#!/usr/bin/env bash

set -e

source build.sh

git tag "$VERSION"
git push origin "$VERSION"

docker build . -t "$DOCKER_TAG"  --build-arg VERSION="$VERSION"

echo "Publishing kdump $VERSION"

docker push "$DOCKER_TAG"
