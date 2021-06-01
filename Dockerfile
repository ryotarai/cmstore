# syntax = docker/dockerfile:experimental
FROM golang:1.16 as builder

WORKDIR /go/src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg go mod download

COPY ./ ./
RUN --mount=type=cache,target=/go/pkg make build

FROM ubuntu:20.04

COPY --from=builder /go/src/bin/cmstore /usr/local/bin/cmstore

ENTRYPOINT ["/usr/local/bin/cmstore"]
