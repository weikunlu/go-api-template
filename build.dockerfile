FROM golang:1.12 as build

ENV GO111MODULE on

WORKDIR /go/release
