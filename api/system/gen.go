package system

//go:generate oapi-codegen --package=system --generate=types,skip-prune  -o ./types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=system --generate=server,skip-prune -o ./server.generated.go openapi.yaml
