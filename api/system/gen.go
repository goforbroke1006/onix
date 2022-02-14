package dashboard_main

//go:generate oapi-codegen --package=system --generate=types,skip-prune  -o ./../../internal/component/api/system/types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=system --generate=server,skip-prune -o ./../../internal/component/api/system/server.generated.go openapi.yaml
