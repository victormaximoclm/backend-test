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

// CreatePartInput representa os dados necessários para criar uma peça. ID e timestamps são gerados pelo service e não usuário
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
