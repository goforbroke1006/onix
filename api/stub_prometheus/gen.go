package stubprometheus

//go:generate oapi-codegen --package=stubprometheus --generate=types,skip-prune  -o ./types.gen.go  openapi.yaml
//go:generate oapi-codegen --package=stubprometheus --generate=server,skip-prune -o ./server.gen.go openapi.yaml
