# go-api-template
startup api template in go

## Prerequisite
- GO v1.12
- PostgreSQL
- Redis

## Setting up the application
- Create a environment file for your application
```
cp .env.example .env
``` 

## Build the project
Build application directly, and find the binary at `./build` 
```
make build
```

Build application without go installation
```
make docker-img
make build-macos-docker
```

Or, type `make help` to check variety actions
```
  all                  clean and run
  build                go build from local
  build-linux          go build for linux
  build-macos          go build for linux
  docker-img           download build image
  build-macos-docker   docker build
  build-linux-docker   docker build for linux
  test                 go test
  coverage             show coverage report in html
  clean                go clean
  run                  build and run go binary
  run-dev              run with live reload
```

## Develop the project

### Install live-reload package for go app
- https://github.com/codegangsta/gin
- https://github.com/cosmtrek/air