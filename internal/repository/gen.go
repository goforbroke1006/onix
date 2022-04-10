package repository

//go:generate mockgen -source=../../domain/service.go     -destination=./mocks/service.go     -package=mocks
//go:generate mockgen -source=../../domain/release.go     -destination=./mocks/release.go     -package=mocks
//go:generate mockgen -source=../../domain/source.go      -destination=./mocks/source.go      -package=mocks
//go:generate mockgen -source=../../domain/criteria.go    -destination=./mocks/criteria.go    -package=mocks
//go:generate mockgen -source=../../domain/measurement.go -destination=./mocks/measurement.go -package=mocks
