package repository

//go:generate mockgen -source=./../../domain/service.go  -destination=./../../mocks/repository/service.generated.go  -package=repository
//go:generate mockgen -source=./../../domain/release.go  -destination=./../../mocks/repository/release.generated.go  -package=repository
//go:generate mockgen -source=./../../domain/source.go   -destination=./../../mocks/repository/source.generated.go   -package=repository
//go:generate mockgen -source=./../../domain/criteria.go -destination=./../../mocks/repository/criteria.generated.go -package=repository
