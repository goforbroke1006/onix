package dashboard_main

//go:generate oapi-codegen --package=dashboard_admin --generate=types,skip-prune  -o ./../../internal/component/api/dashboard_admin/types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=dashboard_admin --generate=server,skip-prune -o ./../../internal/component/api/dashboard_admin/server.generated.go openapi.yaml
