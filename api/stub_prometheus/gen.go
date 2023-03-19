package stubprometheus

//go:generate go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.1
//go:generate mkdir -p ./../../internal/component/stub/prometheus/spec/
//go:generate oapi-codegen --package=spec --generate=types,skip-prune  -o ./../../internal/component/stub/prometheus/spec/types.gen.go  openapi.yaml
//go:generate oapi-codegen --package=spec --generate=server,skip-prune -o ./../../internal/component/stub/prometheus/spec/server.gen.go openapi.yaml
