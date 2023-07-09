#!/bin/sh

GOVERSION=${GOVERSION:-1.20.5}
CONTAINER=${PREFIX}-${ARCH}

# prepare
podman build -t ${CONTAINER} -f - <<EOF
  FROM debian:stable-slim

  RUN dpkg --add-architecture ${ARCH}
  RUN apt-get update && apt-get install --yes --no-install-recommends make curl ca-certificates libgstreamer1.0-dev:${ARCH} libgstreamer-plugins-base1.0-dev:${ARCH} ${DEPS}

  WORKDIR /go
  RUN curl -L "https://go.dev/dl/go${GOVERSION}.linux-amd64.tar.gz" | tar xzf - --strip-components 1
  ENV PATH "$PATH:/go/bin"

  ENV CGO_ENABLED 1
  ENV GOOS "linux"
  ENV GOARCH "${GOARCH}"
  ENV CC "${PREFIX:+$PREFIX-}gcc"
  ENV PKG_CONFIG_PATH "/usr/lib/${PREFIX}/pkgconfig"
EOF

podman run \
  --rm \
  --volume ./:/data \
  --workdir /data \
  ${CONTAINER} \
  go build -mod=vendor -tags=embed -o service-linux-${GOARCH} main.go
