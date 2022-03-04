package repository

//go:generate mockgen -source=./../../domain/service.go -destination=./mocks/service.generated.go -package=mocks
//go:generate mockgen -source=./../../domain/release.go -destination=./mocks/release.generated.go -package=mocks
//go:generate mockgen -source=./../../domain/source.go  -destination=./mocks/source.generated.go -package=mocks
//go:generate mockgen -source=./../../domain/criteria.go  -destination=./mocks/criteria.generated.go -package=mocks
