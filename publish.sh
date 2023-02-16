#!/usr/bin/env bash

set -e

source build.sh
export DOCKER_TAG=gigurra/kdump:$VERSION

echo "Publishing kdump $VERSION"

git tag "$VERSION"
git push origin "$VERSION"

docker build . -t "$DOCKER_TAG"  --build-arg VERSION="$VERSION"
docker push "$DOCKER_TAG"
