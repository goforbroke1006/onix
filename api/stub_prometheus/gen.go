package stubprometheus

//go:generate oapi-codegen --package=stubprometheus --generate=types,skip-prune  -o ./types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=stubprometheus --generate=server,skip-prune -o ./server.generated.go openapi.yaml
