package domain

import "context"

// Service keeps service name.
type Service struct {
	ID string
}

// ServiceRepository describes how to manage Service in db.
type ServiceRepository interface {
	Store(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]Service, error)
}
