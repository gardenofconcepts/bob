.PHONY: clean test run make

DIR = $(shell pwd)

build:
	cd "$(DIR)"
	go build -o bin/builder *.go

run:
	cd "$(DIR)"
	go build -o bin/builder *.go
	./bin/builder -path assets

test:
	go test *.go

clean:
	rm -rf "$(DIR)/bin/*"

init: main.go
	cd "$(DIR)"
	export GOPATH="$(DIR)"
	go get gopkg.in/yaml.v2
	go get github.com/bradfitz/slice
	go get github.com/bmatcuk/doublestar
	go get github.com/aws/aws-sdk-go

default: build
