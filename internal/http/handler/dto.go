package handler

import (
	"time"

	"backend-test/internal/domain"
	"backend-test/internal/service"
)

// PartRequest é o payload aceito em POST/PUT. Separado de domain.Part
// para que o contrato HTTP possa evoluir sem forçar mudança no domínio.
type PartRequest struct {
	Name              string  `json:"name"`
	Category          string  `json:"category"`
	CurrentStock      int     `json:"currentStock"`
	MinimumStock      int     `json:"minimumStock"`
	AverageDailySales float64 `json:"averageDailySales"`
	LeadTimeDays      int     `json:"leadTimeDays"`
	UnitCost          float64 `json:"unitCost"`
	CriticalityLevel  int     `json:"criticalityLevel"`
}

func (r PartRequest) toCreateInput() service.CreatePartInput {
	return service.CreatePartInput{
		Name:              r.Name,
		Category:          r.Category,
		CurrentStock:      r.CurrentStock,
		MinimumStock:      r.MinimumStock,
		AverageDailySales: r.AverageDailySales,
		LeadTimeDays:      r.LeadTimeDays,
		UnitCost:          r.UnitCost,
		CriticalityLevel:  r.CriticalityLevel,
	}
}

func (r PartRequest) toUpdateInput() service.UpdatePartInput {
	return service.UpdatePartInput{
		Name:              r.Name,
		Category:          r.Category,
		CurrentStock:      r.CurrentStock,
		MinimumStock:      r.MinimumStock,
		AverageDailySales: r.AverageDailySales,
		LeadTimeDays:      r.LeadTimeDays,
		UnitCost:          r.UnitCost,
		CriticalityLevel:  r.CriticalityLevel,
	}
}

// PartResponse é o payload retornado ao cliente.
type PartResponse struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Category          string    `json:"category"`
	CurrentStock      int       `json:"currentStock"`
	MinimumStock      int       `json:"minimumStock"`
	AverageDailySales float64   `json:"averageDailySales"`
	LeadTimeDays      int       `json:"leadTimeDays"`
	UnitCost          float64   `json:"unitCost"`
	CriticalityLevel  int       `json:"criticalityLevel"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func newPartResponse(p domain.Part) PartResponse {
	return PartResponse{
		ID:                p.ID,
		Name:              p.Name,
		Category:          p.Category,
		CurrentStock:      p.CurrentStock,
		MinimumStock:      p.MinimumStock,
		AverageDailySales: p.AverageDailySales,
		LeadTimeDays:      p.LeadTimeDays,
		UnitCost:          p.UnitCost,
		CriticalityLevel:  p.CriticalityLevel,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}

func newPartListResponse(parts []domain.Part) []PartResponse {
	result := make([]PartResponse, 0, len(parts))
	for _, p := range parts {
		result = append(result, newPartResponse(p))
	}
	return result
}

// PriorityItemResponse representa uma peça no resultado de priorização
type PriorityItemResponse struct {
	PartID         string  `json:"partId"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	CurrentStock   int     `json:"currentStock"`
	ProjectedStock float64 `json:"projectedStock"`
	MinimumStock   int     `json:"minimumStock"`
	UrgencyScore   float64 `json:"urgencyScore"`
}

// PrioritiesResponse envelopa a lista, replicando o formato
// {"priorities": [...]} do exemplo do desafio.
type PrioritiesResponse struct {
	Priorities []PriorityItemResponse `json:"priorities"`
}

func newPrioritiesResponse(results []domain.PriorityResult) PrioritiesResponse {
	items := make([]PriorityItemResponse, 0, len(results))
	for _, r := range results {
		items = append(items, PriorityItemResponse{
			PartID:         r.Part.ID,
			Name:           r.Part.Name,
			Category:       r.Part.Category,
			CurrentStock:   r.Part.CurrentStock,
			ProjectedStock: r.ProjectedStock,
			MinimumStock:   r.Part.MinimumStock,
			UrgencyScore:   r.UrgencyScore,
		})
	}
	return PrioritiesResponse{Priorities: items}
}

// ErrorResponse padroniza o formato de erro retornado ao cliente.
type ErrorResponse struct {
	Error string `json:"error"`
}
