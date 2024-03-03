#!/bin/sh

GOVERSION=${GOVERSION:-1.22}
CONTAINER=${PREFIX}-${ARCH}

# prepare
podman build -t ${CONTAINER} -f - <<EOF
  FROM golang:${GOVERSION}

  RUN dpkg --add-architecture ${ARCH}
  RUN apt-get update && apt-get install --yes --no-install-recommends make curl ca-certificates libgstreamer1.0-dev:${ARCH} libgstreamer-plugins-base1.0-dev:${ARCH} ${DEPS}

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
  go build -mod=vendor -tags=embed -ldflags "-s -w" -trimpath -o service-linux-${GOARCH} main.go
