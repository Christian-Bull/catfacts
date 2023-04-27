#! /bin/bash

platarch=linux/amd64,linux/arm64
buildah build --jobs=2 --platform=$platarch --manifest catfacts .

shortSHA=$(git rev-parse --short HEAD)

buildah tag localhost/catfacts csbull55/catfacts:$shortSHA

podman manifest rm localhost/catfacts
podman manifest push --all csbull55/catfacts:$shortSHA docker://csbull55/catfacts:$shortSHA
podman manifest push --all csbull55/catfacts:$shortSHA docker://csbull55/catfacts:$shortSHA-latest
