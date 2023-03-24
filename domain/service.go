package domain

import "context"

// Service keeps service name.
type Service struct {
	ID string
}

// ServiceStorage describes how to manage Service in db.
type ServiceStorage interface {
	Store(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]Service, error)
}
