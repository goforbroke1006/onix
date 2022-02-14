package dashboard_main

//go:generate oapi-codegen --package=dashboard_main --generate=types,skip-prune  -o ./../../internal/component/api/dashboard_main/types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=dashboard_main --generate=server,skip-prune -o ./../../internal/component/api/dashboard_main/server.generated.go openapi.yaml
