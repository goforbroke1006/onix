package stub_prometheus

//go:generate oapi-codegen --package=stub_prometheus --generate=types,skip-prune  -o ./types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=stub_prometheus --generate=server,skip-prune -o ./server.generated.go openapi.yaml
