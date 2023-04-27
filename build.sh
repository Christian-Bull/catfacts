#! /bin/bash

podman buildx build -t csbul55/catfacts:v7 .
podman buildx build -t catfacts:v2 .
podman push csbul55/catfacts:v7


platarch=linux/amd64,linux/arm64
buildah build --jobs=2 --platform=$platarch --manifest catfacts .