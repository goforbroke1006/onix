package dashboard_admin

//go:generate oapi-codegen --package=dashboard_admin --generate=types,skip-prune  -o ./types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=dashboard_admin --generate=server,skip-prune -o ./server.generated.go openapi.yaml
