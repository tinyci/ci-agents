# vim: ft=dockerfile

FROM ubuntu:20.04

ARG GO_VERSION=1.15.8
ARG SWAGGER_VERSION=v0.18.0
ARG PROTOC_VERSION=3.15.6
ARG TIMEZONE=Etc/UTC
ARG APT_MIRROR=mirror.pnl.gov

ENV TZ=${TIMEZONE}
RUN ln -fs /usr/share/zoneinfo/${TIMEZONE} /etc/localtime

RUN perl -i.bak -pe "s!(security|archive)\.ubuntu\.com!${APT_MIRROR}!g" /etc/apt/sources.list
RUN apt-get update && apt-get dist-upgrade -y && apt-get install -y \
  curl \
  wget \
  gnupg \
  git \
  mercurial \
  build-essential \
  sudo \
  libnss3-tools \
  unzip

ENV GOPATH=/go \
    PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

RUN curl -sSL "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" | tar xz -C /usr/local
RUN mkdir /go

RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

ENV protoc_fn=protoc-${PROTOC_VERSION}-linux-x86_64.zip

RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${protoc_fn}
RUN unzip "${protoc_fn}" -d /usr && rm -f "${protoc_fn}"
RUN chmod +x /usr/bin/protoc && chmod -R 755 /usr/include/google

RUN curl -sSL "https://github.com/go-swagger/go-swagger/releases/download/${SWAGGER_VERSION}/swagger_linux_amd64" >/go/bin/swagger && chmod +x /go/bin/swagger
