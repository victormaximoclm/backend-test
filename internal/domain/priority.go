package domain

import (
	"math"
	"sort"
)

// Estrutura da Regra de Negócio, consumo, estoque projetado e score urgência são floats porque para chegar em seus valores
// são utilizados numeros que podem possuir casas decimais.
type PriorityResult struct {
	Part                Part
	ExpectedConsumption float64
	ProjectedStock      float64
	NeedsRestock        bool
	UrgencyScore        float64
}

// função para arredondar para evitar muitas casas decimais
func roundToTwoDecimals(v float64) float64 {
	return math.Round(v*100) / 100
}

// Aplicação das regras de negócio
func CalculatePriority(p Part) PriorityResult {
	expectedConsumption := p.AverageDailySales * float64(p.LeadTimeDays)
	projectedStock := float64(p.CurrentStock) - expectedConsumption

	needsRestock := projectedStock < float64(p.MinimumStock)

	urgencyScore := (float64(p.MinimumStock) - projectedStock) * float64(p.CriticalityLevel)

	return PriorityResult{
		Part:                p,
		ExpectedConsumption: roundToTwoDecimals(expectedConsumption),
		ProjectedStock:      roundToTwoDecimals(projectedStock),
		NeedsRestock:        needsRestock,
		UrgencyScore:        roundToTwoDecimals(urgencyScore),
	}
}

// Percorrer lista de peças e filtra quem precisa de restoque e aplica função de ordenação
func RankPriorities(parts []Part) []PriorityResult {
	results := make([]PriorityResult, 0, len(parts))
	for _, p := range parts {
		r := CalculatePriority(p)
		if r.NeedsRestock {
			results = append(results, r)
		}
	}

	sortPriorityResults(results)
	return results
}

// função de ordenação por prioridade
func sortPriorityResults(results []PriorityResult) {
	sort.Slice(results, func(i, j int) bool {
		a, b := results[i], results[j]

		if a.UrgencyScore != b.UrgencyScore {
			return a.UrgencyScore > b.UrgencyScore
		}

		if a.Part.CriticalityLevel != b.Part.CriticalityLevel {
			return a.Part.CriticalityLevel > b.Part.CriticalityLevel
		}

		if a.Part.AverageDailySales != b.Part.AverageDailySales {
			return a.Part.AverageDailySales > b.Part.AverageDailySales
		}

		return a.Part.Name < b.Part.Name
	})
}
