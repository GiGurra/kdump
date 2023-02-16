#!/usr/bin/env bash

set -e

if [ -n "$(git status --porcelain)" ]; then
  echo "Uncommitted changes detected - Bailing!"
  exit 1
fi

source build.sh
export DOCKER_TAG=gigurra/kdump:$VERSION

echo "Publishing kdump $VERSION"

git tag "$VERSION"
git push origin "$VERSION"

docker build . -t "$DOCKER_TAG" --platform linux/amd64 --build-arg VERSION="$VERSION"
docker push "$DOCKER_TAG"
