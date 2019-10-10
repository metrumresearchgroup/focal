#!/bin/bash

OS=linux
ARCH=amd64
CONTAINER=dukeofubuntu/ugo13

docker run --rm -v $(pwd):/app -e GOOS=${OS} -e GOARCH=${ARCH} ${CONTAINER} go build -o focal *.go