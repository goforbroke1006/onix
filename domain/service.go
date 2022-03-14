package domain

// Service keeps service name.
type Service struct {
	Title string
}

// ServiceRepository describes how to manage Service in db.
type ServiceRepository interface {
	Store(title string) error
	GetAll() ([]Service, error)
}
