.PHONY: build clean test run make cleancode init dist

DIR = $(shell pwd)
REVISION = $(shell git rev-parse --short HEAD)

build:
	cd "$(DIR)"
	sed -i -r 's/^(const APP_BUILD string = )"([a-zA-Z0-9]+)"/\1"$(REVISION)"/' src/bob/version.go
	go build -o bin/bob $(shell find src/bob/*.go | grep -v _test )

dist:
	# GOARCH=386 = 32bit
	sed -i -r 's/^(const APP_BUILD string = )"([a-zA-Z0-9]+)"/\1"$(REVISION)"/' src/bob/version.go
	env GOOS=linux   GOARCH=amd64 go build -o bin/linux/amd64/bob src/bob/*.go
	env GOOS=darwin  GOARCH=amd64 go build -o bin/darwin/amd64/bob src/bob/*.go
	env GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/bob src/bob/*.go
	tar -C bin/linux/amd64 -cvzf build/bob_linux_amd64.tar.gz bob
	tar -C bin/darwin/amd64 -cvzf build/bob_darwin_amd64.tar.gz bob
	tar -C bin/windows/amd64 -cvzf build/bob_windows_amd64.tar.gz bob

run:
	cd "$(DIR)"
	go build -o bin/bob src/bob/*.go
	./bin/bob -path assets

test:
	go test -v src/bob/analyzer/*.go
	go test -v src/bob/archive/*.go
	go test -v src/bob/builder/*.go
	go test -v src/bob/config/*.go
	go test -v src/bob/hash/*.go
	go test -v src/bob/parser/*.go
	go test -v src/bob/path/*.go
	go test -v src/bob/reader/*.go
	go test -v src/bob/storage/*.go
	go test -v src/bob/util/*.go
	go test -v src/bob/*.go

test-coverage:
	go test -coverprofile=coverage.out

	cat coverage.out | cut -c30-

	sed -i "s/_\/home\/dennis\/htdocs\/builder/./" coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm -rf coverage.out

clean:
	rm -rf "$(DIR)/bin/*" ; \
	rm -rf "$(DIR)/build/*"

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
	go get github.com/Knetic/govaluate
	go get github.com/hashicorp/go-version
	go get github.com/patrickmn/go-cache

cleancode:
	cd "$(DIR)"
	gofmt -w src/bob/*.go

default: build
