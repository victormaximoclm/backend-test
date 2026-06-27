package service

import (
	"backend-test/internal/domain"
	"backend-test/internal/repository"

	"github.com/google/uuid"
)

// Part Service é onde vai ter os casos de uso de CRUD de peças, dependendo só da interface repository.PartRepository
type PartService struct {
	repo repository.PartRepository
}

// Construtor Cria um PartService e associa o repositório a ele.
func NewPartService(repo repository.PartRepository) *PartService {
	return &PartService{repo: repo}
}

// CreatePartInput representa os dados necessários para criar uma peça. ID e timestamps são gerados pelo service e repositório respectivamente, e não usuário
type CreatePartInput struct {
	Name              string
	Category          string
	CurrentStock      int
	MinimumStock      int
	AverageDailySales float64
	LeadTimeDays      int
	UnitCost          float64
	CriticalityLevel  int
}

// CreatePart cria uma nova peça, valida os dados e delega a persistência ao repositório.
func (s *PartService) CreatePart(input CreatePartInput) (domain.Part, error) {
	part := domain.Part{
		ID:                uuid.NewString(),
		Name:              input.Name,
		Category:          input.Category,
		CurrentStock:      input.CurrentStock,
		MinimumStock:      input.MinimumStock,
		AverageDailySales: input.AverageDailySales,
		LeadTimeDays:      input.LeadTimeDays,
		UnitCost:          input.UnitCost,
		CriticalityLevel:  input.CriticalityLevel,
	}

	if err := part.Validate(); err != nil {
		return domain.Part{}, err
	}

	return s.repo.Create(part)
}

// GetPart procura peça usando método de busca do repositório
func (s *PartService) GetPart(id string) (domain.Part, error) {
	return s.repo.GetByID(id)
}

// ListParts usa método de listagem do repositório
func (s *PartService) ListParts() ([]domain.Part, error) {
	return s.repo.List()
}

// ListPartsByCategory usa método listagem por categoria do repositório
func (s *PartService) ListPartsByCategory(category string) ([]domain.Part, error) {
	return s.repo.ListByCategory(category)
}

// DeletePart usa método delete do repositório
func (s *PartService) DeletePart(id string) error {
	return s.repo.Delete(id)
}

// Estrutura utilizada para atualização de peça
type UpdatePartInput struct {
	Name              string
	Category          string
	CurrentStock      int
	MinimumStock      int
	AverageDailySales float64
	LeadTimeDays      int
	UnitCost          float64
	CriticalityLevel  int
}

// UpdatePart busca a peça, atualiza os dados, valida os dados e delega a persistência ao repositório.
func (s *PartService) UpdatePart(id string, input UpdatePartInput) (domain.Part, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return domain.Part{}, err
	}

	updated := domain.Part{
		ID:                existing.ID,
		Name:              input.Name,
		Category:          input.Category,
		CurrentStock:      input.CurrentStock,
		MinimumStock:      input.MinimumStock,
		AverageDailySales: input.AverageDailySales,
		LeadTimeDays:      input.LeadTimeDays,
		UnitCost:          input.UnitCost,
		CriticalityLevel:  input.CriticalityLevel,
	}

	if err := updated.Validate(); err != nil {
		return domain.Part{}, err
	}

	return s.repo.Update(updated)
}
