package domain

type Service struct {
	Title string
}

type ServiceRepository interface {
	Story(title string) error
	GetAll() ([]Service, error)
}
