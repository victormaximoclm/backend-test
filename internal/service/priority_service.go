package service

import (
	"backend-test/internal/domain"
	"backend-test/internal/repository"
)

// PriorityService aplica regras de negócio e define a ordem de prioridade das peças.
type PriorityService struct {
	repo repository.PartRepository
}

func NewPriorityService(repo repository.PartRepository) *PriorityService {
	return &PriorityService{repo: repo}
}

// GetPriorities retorna as peças que precisam de reposição, já ordenadas
// por urgência conforme as regras de negócio do domínio.
func (s *PriorityService) GetPriorities() ([]domain.PriorityResult, error) {
	parts, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	return domain.RankPriorities(parts), nil
}
