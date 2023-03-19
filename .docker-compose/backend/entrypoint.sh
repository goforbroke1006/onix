#!/bin/bash

go mod download -x
go generate ./...
find ./ -name "*.gen.go" -exec chmod 0777 {} \;
find ./ -name "*.mock.go" -exec chmod 0777 {} \;

appArgs="$*"

echo "Run application with args '$appArgs'"

CompileDaemon --build='go build -o /tmp/application .' --command="/tmp/application $appArgs"
