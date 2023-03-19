package v1

//go:generate go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.1
//go:generate mkdir -p ./../../../internal/component/api/external/v1/spec/
//go:generate oapi-codegen --package=spec --generate=types,skip-prune  -o ./../../../internal/component/api/external/v1/spec/types.gen.go  openapi.yaml
//go:generate oapi-codegen --package=spec --generate=server,skip-prune -o ./../../../internal/component/api/external/v1/spec/server.gen.go openapi.yaml
