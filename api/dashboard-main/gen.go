package dashboard_main

//go:generate oapi-codegen --package=dashboard_main --generate=types,skip-prune  -o ./types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=dashboard_main --generate=server,skip-prune -o ./server.generated.go openapi.yaml
