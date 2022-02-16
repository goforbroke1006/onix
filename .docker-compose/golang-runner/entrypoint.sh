#!/bin/bash

go get -u github.com/deepmap/oapi-codegen/cmd/oapi-codegen

go mod download

go generate ./...
find ./ -name '*.generated.go' | xargs chmod 0777 '{}'

# shellcheck disable=SC2068
go run ./main.go $@
