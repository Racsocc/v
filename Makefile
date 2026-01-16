VERSION := $(shell git describe --tags --dirty --always)

install:
	go install -ldflags "-X main.Version=$(VERSION)" .

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o v .
