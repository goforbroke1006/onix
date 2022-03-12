package dashboardmain

//go:generate oapi-codegen --package=dashboardmain --generate=types,skip-prune  -o ./types.generated.go  openapi.yaml
//go:generate oapi-codegen --package=dashboardmain --generate=server,skip-prune -o ./server.generated.go openapi.yaml
