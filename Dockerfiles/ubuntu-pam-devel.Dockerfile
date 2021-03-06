FROM ubuntu

WORKDIR /tmp

#Prep Go 13
ENV GOROOT=/usr/local/go
ENV GOPATH=/root/go
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin

RUN     apt-get update && \
        apt-get install -y git wget build-essential libpam0g-dev curl  && \
        wget https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz && \
        tar -C /usr/local -xvzf go1.13.1.linux-amd64.tar.gz && \
        curl -L https://github.com/goreleaser/goreleaser/releases/latest/download/goreleaser_amd64.deb > goreleaser.deb && \
        dpkg -i goreleaser.deb

WORKDIR /app