# vim: ft=dockerfile
FROM ubuntu:20.04

ARG GO_VERSION=1.16.3
ARG POSTGRES_VERSION=11
ARG PROTOC_VERSION=3.15.6
ARG CADDY_VERSION=2.3.0
ARG TIMEZONE=Etc/UTC
ARG APT_MIRROR=mirror.pnl.gov
ARG TESTING=

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

RUN curl -sSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN echo 'deb http://apt.postgresql.org/pub/repos/apt/ focal-pgdg main' | tee -a /etc/apt/sources.list.d/postgresql.list

ENV GOPATH=/go \
    TESTING=${TESTING} \
    CAROOT=/var/ca \
    PATH=/usr/lib/postgresql/${POSTGRES_VERSION}/bin:/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

RUN apt-get update -qq && apt-get install "postgresql-${POSTGRES_VERSION}" -y -qq

RUN curl -sSL "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" | tar xz -C /usr/local
RUN mkdir /go

RUN go get github.com/erikh/migrator
RUN go get filippo.io/mkcert
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
RUN go get -u github.com/deepmap/oapi-codegen/cmd/oapi-codegen

RUN wget https://github.com/caddyserver/caddy/releases/download/v2.3.0/caddy_${CADDY_VERSION}_linux_amd64.deb
RUN dpkg -i caddy_${CADDY_VERSION}_linux_amd64.deb

ENV protoc_fn=protoc-${PROTOC_VERSION}-linux-x86_64.zip

RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${protoc_fn}
RUN unzip "${protoc_fn}" -d /usr && rm -f "${protoc_fn}" && chmod -R 755 /usr/bin/protoc /usr/include/google

COPY .config/Caddyfile /Caddyfile
COPY entrypoint.sh /entrypoint.sh
RUN chmod 755 /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
CMD []
