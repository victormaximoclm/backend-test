package service

import (
	"backend-test/internal/domain"
	"backend-test/internal/repository"
)

// Fake do PartRepository usado em testes para evitar dependência de banco de dados e infraestrutura
type fakePartRepository struct {
	parts map[string]domain.Part
}

func newFakePartRepository() *fakePartRepository {
	return &fakePartRepository{parts: make(map[string]domain.Part)}
}

func (f *fakePartRepository) Create(part domain.Part) (domain.Part, error) {
	f.parts[part.ID] = part
	return part, nil
}

func (f *fakePartRepository) GetByID(id string) (domain.Part, error) {
	p, exists := f.parts[id]
	if !exists {
		return domain.Part{}, repository.ErrPartNotFound
	}
	return p, nil
}

func (f *fakePartRepository) List() ([]domain.Part, error) {
	result := make([]domain.Part, 0, len(f.parts))
	for _, p := range f.parts {
		result = append(result, p)
	}
	return result, nil
}

func (f *fakePartRepository) ListByCategory(category string) ([]domain.Part, error) {
	result := make([]domain.Part, 0)
	for _, p := range f.parts {
		if p.Category == category {
			result = append(result, p)
		}
	}
	return result, nil
}

func (f *fakePartRepository) Update(part domain.Part) (domain.Part, error) {
	if _, exists := f.parts[part.ID]; !exists {
		return domain.Part{}, repository.ErrPartNotFound
	}
	f.parts[part.ID] = part
	return part, nil
}

func (f *fakePartRepository) Delete(id string) error {
	if _, exists := f.parts[id]; !exists {
		return repository.ErrPartNotFound
	}
	delete(f.parts, id)
	return nil
}
