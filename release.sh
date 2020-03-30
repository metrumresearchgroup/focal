#!/bin/bash

IMAGE=dukeofubuntu/goreleaser
DIRECTORY=/root/go/github.com/metrumresearchgroup/focal

#docker run --rm -it -v $(pwd):$DIRECTORY -w $DIRECTORY/cmd/focal -e GITHUB_TOKEN=$GITHUB_TOKEN dukeofubuntu/goreleaser goreleaser --skip-validate --skip-publish --rm-dist --snapshot

#Githubtoken should come from env of whoever is running
docker run --rm -it -v $(pwd):$DIRECTORY -w $DIRECTORY -e GITHUB_TOKEN=$GITHUB_TOKEN dukeofubuntu/goreleaser goreleaser --rm-dist

