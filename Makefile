
build:
	CGO_ENABLED=1 go build -mod=vendor -o service main.go

build-static:
	CGO_ENABLED=1 go build -mod=vendor -tags=embed -o service main.go