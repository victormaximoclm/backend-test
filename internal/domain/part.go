package domain

import (
	"errors"
	"time"
)

// Estrutura do estoque
type Part struct {
	ID                string
	Name              string
	Category          string
	CurrentStock      int
	MinimumStock      int
	AverageDailySales float64
	LeadTimeDays      int
	UnitCost          float64
	CriticalityLevel  int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

const (
	MinCriticalityLevel = 1
	MaxCriticalityLevel = 5
)

var (
	ErrEmptyName            = errors.New("name não pode ser vazio")
	ErrEmptyCategory        = errors.New("category não pode ser vazio")
	ErrNegativeMinimumStock = errors.New("minimumStock não pode ser negativo")
	ErrNegativeDailySales   = errors.New("averageDailySales não pode ser negativo")
	ErrNegativeLeadTime     = errors.New("leadTimeDays não pode ser negativo")
	ErrNegativeUnitCost     = errors.New("unitCost não pode ser negativo")
	ErrInvalidCriticality   = errors.New("criticalityLevel deve estar entre 1 e 5")
)

// Validações mínimas da estrutura
func (p *Part) Validate() error {
	if p.Name == "" {
		return ErrEmptyName
	}

	if p.Category == "" {
		return ErrEmptyCategory
	}

	if p.MinimumStock < 0 {
		return ErrNegativeMinimumStock
	}

	if p.AverageDailySales < 0 {
		return ErrNegativeDailySales
	}

	if p.LeadTimeDays < 0 {
		return ErrNegativeLeadTime
	}

	if p.UnitCost < 0 {
		return ErrNegativeUnitCost
	}

	if p.CriticalityLevel < MinCriticalityLevel || p.CriticalityLevel > MaxCriticalityLevel {
		return ErrInvalidCriticality
	}
	return nil
}
