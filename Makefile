.PHONY: build clean test run make cleancode init dist

DIR = $(shell pwd)
REVISION = $(shell git rev-parse --short HEAD)

build:
	cd "$(DIR)"
	go build -o bin/builder *.go

dist:
	# GOARCH=386 = 32bit
	sed -i -r 's/^(const APP_BUILD string = )"([a-zA-Z0-9]+)"/\1"$(REVISION)"/' version.go
	env GOOS=linux   GOARCH=amd64 go build -o bin/linux/amd64/builder *.go
	env GOOS=darwin  GOARCH=amd64 go build -o bin/darwin/amd64/builder *.go
	env GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/builder *.go
	tar -C bin/linux/amd64 -cvzf build/builder_linux_amd64.tar.gz builder
	tar -C bin/darwin/amd64 -cvzf build/builder_darwin_amd64.tar.gz builder
	tar -C bin/windows/amd64 -cvzf build/builder_windows_amd64.tar.gz builder

run:
	cd "$(DIR)"
	go build -o bin/builder *.go
	./bin/builder -path assets

test:
	go test *.go

clean:
	rm -rf "$(DIR)/bin/*"

init:
	cd "$(DIR)"
	export GOPATH="$(DIR)"
	go get gopkg.in/yaml.v2
	go get gopkg.in/urfave/cli.v1
	go get github.com/bradfitz/slice
	go get github.com/bmatcuk/doublestar
	go get github.com/aws/aws-sdk-go
	go get github.com/Sirupsen/logrus
	go get github.com/imdario/mergo

cleancode:
	cd "$(DIR)"
	gofmt -w *.go

default: build
