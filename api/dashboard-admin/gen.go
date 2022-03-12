package dashboardadmin

//go:generate oapi-codegen --package=dashboardadmin --generate=types,skip-prune  -o ./types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=dashboardadmin --generate=server,skip-prune -o ./server.generated.go openapi.yaml
