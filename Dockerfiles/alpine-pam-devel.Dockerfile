FROM golang:1.13-alpine3.10

WORKDIR /app
RUN apk add build-base linux-pam-dev