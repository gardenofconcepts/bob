.PHONY: build clean test run make cleancode init dist

DIR = $(shell pwd)
REVISION = $(shell git rev-parse --short HEAD)

build:
	cd "$(DIR)"
	go build -o bin/bob $(shell find *.go | grep -v _test | grep -v s3 )

dist:
	# GOARCH=386 = 32bit
	sed -i -r 's/^(const APP_BUILD string = )"([a-zA-Z0-9]+)"/\1"$(REVISION)"/' version.go
	env GOOS=linux   GOARCH=amd64 go build -o bin/linux/amd64/bob *.go
	env GOOS=darwin  GOARCH=amd64 go build -o bin/darwin/amd64/bob *.go
	env GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/bob *.go
	tar -C bin/linux/amd64 -cvzf build/bob_linux_amd64.tar.gz bob
	tar -C bin/darwin/amd64 -cvzf build/bob_darwin_amd64.tar.gz bob
	tar -C bin/windows/amd64 -cvzf build/bob_windows_amd64.tar.gz bob

run:
	cd "$(DIR)"
	go build -o bin/bob *.go
	./bin/bob -path assets

test:
	go test -v *.go

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

cleancode:
	cd "$(DIR)"
	gofmt -w *.go

default: build
