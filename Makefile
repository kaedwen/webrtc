
build:
	go build -mod=vendor -o service main.go

build-static:
	go build -mod=vendor -tags=embed -o service main.go