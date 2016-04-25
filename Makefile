.PHONY: clean

DIR = $(shell pwd)

make: main.go
	cd "$(DIR)"
	go build -o bin/builder main.go

run: main.go
	cd "$(DIR)"
	go run main.go -path assets

clean:
	rm -rf "$(DIR)/bin/*"

init: main.go
	cd "$(DIR)"
	export GOPATH="$(DIR)"
	go get gopkg.in/yaml.v2
	go get github.com/bradfitz/slice
	go get github.com/bmatcuk/doublestar
