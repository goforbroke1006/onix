package prometheus

//go:generate oapi-codegen --package=prometheus --generate=types,skip-prune  -o ./types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=prometheus --generate=server,skip-prune -o ./server.generated.go openapi.yaml
