# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
BINARY_NAME=api

# Binary parameters
VERSION=$(shell git describe --tags)
DATE=`date +%FT%T%z`
BUILD=$(shell git rev-parse HEAD)
CURRDIR = $(shell pwd)

LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION} -X main.BuildDate=${DATE} -X main.Build=$(BUILD)"

.PHONY: all build test clean

## all: clean and run
all: clean run

## build: go build from local
build:
	$(GOBUILD) ${LDFLAGS} -o ./build/$(BINARY_NAME) -tags=jsoniter -v ./

## build-linux: go build for linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) ${LDFLAGS} -o ./build/$(BINARY_NAME) -tags=jsoniter -v ./

## build-macos: go build for linux
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) ${LDFLAGS} -o ./build/$(BINARY_NAME) -tags=jsoniter -v ./

## docker-img: download build image
docker-img:
	docker build -t kkday/golang:v1.12 -f ./build.dockerfile .

## build-macos-docker: docker build
build-macos-docker:
	docker run -v ${CURRDIR}:/go/release:rw --rm kkday/golang:v1.12 bash -c "make build-darwin"

## build-linux-docker: docker build for linux
build-linux-docker:
	docker run -v ${CURRDIR}:/go/release:rw --rm kkday/golang:v1.12 bash -c "make build-linux"

## test: go test
test:
	TESTING=true $(GOTEST) -coverprofile=code_coverage.out -v ./...

## coverage: show coverage report in html
coverage:
	$(GOTOOL) cover -html=code_coverage.out

## clean: go clean
clean:
	$(GOCLEAN)
	rm -f ./build/$(BINARY_NAME)
	rm ./logs/*
	rm code_coverage.*

## run: build and run go binary
run:
	$(GOBUILD) ${LDFLAGS} -o ./build/$(BINARY_NAME) -tags=jsoniter -v ./
	./build/$(BINARY_NAME) runserver

# https://github.com/codegangsta/gin
# https://github.com/cosmtrek/air
## run-dev: run with live reload
run-dev:
	DEVELOP=true ~/go/bin/gin --port 3000 --appPort 3001 main.go

help: Makefile
	@echo
	@echo " Choose a command run in "$(BINARY_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
