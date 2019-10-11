#!/bin/bash

IMAGE=dukeofubuntu/goreleaser
DIRECTORY=/root/go/github.com/metrumresearchgroup/focal

#Githubtoken should come from env of whoever is running
docker run --rm -it -v $(pwd):$DIRECTORY -w $DIRECTORY -e GITHUB_TOKEN=$GITHUB_TOKEN dukeofubuntu/goreleaser goreleaser --snapshot --skip-publish --rm-dist

