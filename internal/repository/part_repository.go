package repository

import "backend-test/internal/domain"

//interface obrigatória que repositórios devem ter
type PartRepository interface {
	Create(part domain.Part) (domain.Part, error)
	GetByID(id string) (domain.Part, error)
	List() ([]domain.Part, error)
	ListByCategory(category string) ([]domain.Part, error)
	Update(part domain.Part) (domain.Part, error)
	Delete(id string) error
}
