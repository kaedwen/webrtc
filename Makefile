build:
	CGO_ENABLED=1 go build -mod=vendor -tags=embed -o service main.go

build-static:
	npm --prefix static ci && npm --prefix static run build

build-armhf:
	GOARCH=arm \
	ARCH=armhf \
	DEPS="gcc-arm-linux-gnueabihf libc6-dev-armhf-cross" \
	PREFIX=arm-linux-gnueabihf \
	./build-cross.sh

build-arm64:
	GOARCH=arm64 \
	ARCH=arm64 \
	DEPS="gcc-aarch64-linux-gnu libc6-dev-arm64-cross" \
	PREFIX=aarch64-linux-gnu \
	./build-cross.sh
