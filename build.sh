#! /bin/bash

podman buildx build -t csbul55/catfacts:v7 .
podman buildx build -t catfacts:v2 .
podman push csbul55/catfacts:v7


platarch=linux/amd64,linux/arm64
buildah build --jobs=2 --platform=$platarch --manifest catfacts .
buildah tag localhost/catfacts csbull55/catfacts:v10
podman manifest rm localhost/catfacts
podman manifest push --all csbull55/catfacts:v10 docker://csbull55/catfacts:v10
