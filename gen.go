package main

//go:generate mockgen -source=./domain/service.go  -destination=./mocks/repository/service.go  -package=repository
//go:generate mockgen -source=./domain/release.go  -destination=./mocks/repository/release.go  -package=repository
//go:generate mockgen -source=./domain/source.go   -destination=./mocks/repository/source.go   -package=repository
//go:generate mockgen -source=./domain/criteria.go -destination=./mocks/repository/criteria.go -package=repository
