package domain

import "testing"

// newTestPart cria uma peça base válida, e permite sobrescrever só os
// campos que importam para cada cenário de teste
func newTestPart(overrides func(*Part)) Part {
	p := Part{
		ID:                "test-id",
		Name:              "Peça Teste",
		Category:          "engine",
		CurrentStock:      15,
		MinimumStock:      20,
		AverageDailySales: 4,
		LeadTimeDays:      5,
		UnitCost:          18.50,
		CriticalityLevel:  3,
	}
	if overrides != nil {
		overrides(&p)
	}
	return p
}

func TestCalculatePriority(t *testing.T) {
	// expectedConsumption = 4 * 5 = 20
	// projectedStock = 15 - 20 = -5
	// needsRestock = -5 < 20 -> true
	// urgencyScore = (20 - (-5)) * 3 = 75
	p := newTestPart(nil)

	result := CalculatePriority(p)

	if result.ExpectedConsumption != 20 {
		t.Errorf("expectedConsumption = %v; esperado 20", result.ExpectedConsumption)
	}
	if result.ProjectedStock != -5 {
		t.Errorf("projectedStock = %v; esperado -5", result.ProjectedStock)
	}
	if !result.NeedsRestock {
		t.Errorf("needsRestock = false; esperado true")
	}
	if result.UrgencyScore != 75 {
		t.Errorf("urgencyScore = %v; esperado 75", result.UrgencyScore)
	}
}

func TestCalculatePriority_NaoPrecisaReposicao(t *testing.T) {
	// Estoque alto e venda baixa: projectedStock fica >= minimumStock.
	p := newTestPart(func(p *Part) {
		p.CurrentStock = 100
		p.MinimumStock = 20
		p.AverageDailySales = 1
		p.LeadTimeDays = 2
	})

	result := CalculatePriority(p)

	// expectedConsumption = 2, projectedStock = 98, 98 < 20? não.
	if result.NeedsRestock {
		t.Errorf("needsRestock = true; esperado false (estoque suficiente)")
	}
}

func TestCalculatePriority_EstoqueAtualNegativo(t *testing.T) {
	// Cenário extremo pedido explicitamente no desafio.
	p := newTestPart(func(p *Part) {
		p.CurrentStock = -10
		p.MinimumStock = 20
		p.AverageDailySales = 2
		p.LeadTimeDays = 3
		p.CriticalityLevel = 5
	})

	result := CalculatePriority(p)

	// expectedConsumption = 6, projectedStock = -10 - 6 = -16
	if result.ProjectedStock != -16 {
		t.Errorf("projectedStock = %v; esperado -16", result.ProjectedStock)
	}
	if !result.NeedsRestock {
		t.Errorf("peça com estoque negativo deveria sempre precisar de reposição")
	}
	// urgencyScore = (20 - (-16)) * 5 = 180
	if result.UrgencyScore != 180 {
		t.Errorf("urgencyScore = %v; esperado 180", result.UrgencyScore)
	}
}

func TestCalculatePriority_VendaMediaZero(t *testing.T) {
	// averageDailySales = 0: expectedConsumption deve ser 0.
	p := newTestPart(func(p *Part) {
		p.CurrentStock = 10
		p.MinimumStock = 20
		p.AverageDailySales = 0
		p.LeadTimeDays = 30 // lead time alto, mas irrelevante aqui
	})

	result := CalculatePriority(p)

	if result.ExpectedConsumption != 0 {
		t.Errorf("expectedConsumption = %v; esperado 0", result.ExpectedConsumption)
	}
	if result.ProjectedStock != 10 {
		t.Errorf("projectedStock = %v; esperado 10 (igual ao currentStock)", result.ProjectedStock)
	}
	if !result.NeedsRestock {
		t.Errorf("needsRestock = false; esperado true (estoque já abaixo do mínimo)")
	}
}

func TestCalculatePriority_VendaZeroEEstoqueSuficiente(t *testing.T) {
	p := newTestPart(func(p *Part) {
		p.CurrentStock = 25
		p.MinimumStock = 20
		p.AverageDailySales = 0
		p.LeadTimeDays = 10
	})

	result := CalculatePriority(p)

	if result.NeedsRestock {
		t.Errorf("needsRestock = true; esperado false")
	}
}

func TestCalculatePriority_LeadTimeAlto(t *testing.T) {
	p := newTestPart(func(p *Part) {
		p.CurrentStock = 50
		p.MinimumStock = 20
		p.AverageDailySales = 10
		p.LeadTimeDays = 90
		p.CriticalityLevel = 2
	})

	result := CalculatePriority(p)

	// expectedConsumption = 900, projectedStock = 50 - 900 = -850
	if result.ExpectedConsumption != 900 {
		t.Errorf("expectedConsumption = %v; esperado 900", result.ExpectedConsumption)
	}
	if result.ProjectedStock != -850 {
		t.Errorf("projectedStock = %v; esperado -850", result.ProjectedStock)
	}
	// urgencyScore = (20 - (-850)) * 2 = 1740
	if result.UrgencyScore != 1740 {
		t.Errorf("urgencyScore = %v; esperado 1740", result.UrgencyScore)
	}
}

func TestRankPriorities_FiltraApenasQuemPrecisaReposicao(t *testing.T) {
	parts := []Part{
		newTestPart(func(p *Part) {
			p.ID = "ok"
			p.Name = "Peça OK"
			p.CurrentStock = 1000
			p.MinimumStock = 10
			p.AverageDailySales = 1
			p.LeadTimeDays = 1
		}),
		newTestPart(func(p *Part) {
			p.ID = "urgente"
			p.Name = "Peça Urgente"
			p.CurrentStock = 5
			p.MinimumStock = 20
			p.AverageDailySales = 4
			p.LeadTimeDays = 5
		}),
	}

	results := RankPriorities(parts)

	if len(results) != 1 {
		t.Fatalf("esperado 1 resultado, obtido %d", len(results))
	}
	if results[0].Part.ID != "urgente" {
		t.Errorf("esperado peça 'urgente', obtido %s", results[0].Part.ID)
	}
}

func TestRankPriorities_DesempatePorCriticalityLevel(t *testing.T) {
	// Peça A: projectedStock=10, diff=10, criticality=2 -> score=20
	// Peça B: projectedStock=15, diff=5,  criticality=4 -> score=20
	a := newTestPart(func(p *Part) {
		p.ID = "a"
		p.Name = "A"
		p.CurrentStock = 10
		p.MinimumStock = 20
		p.AverageDailySales = 0
		p.LeadTimeDays = 0
		p.CriticalityLevel = 2
	})
	b := newTestPart(func(p *Part) {
		p.ID = "b"
		p.Name = "B"
		p.CurrentStock = 15
		p.MinimumStock = 20
		p.AverageDailySales = 0
		p.LeadTimeDays = 0
		p.CriticalityLevel = 4
	})

	results := RankPriorities([]Part{a, b})

	if results[0].UrgencyScore != results[1].UrgencyScore {
		t.Fatalf("pré-condição falhou: scores deveriam ser iguais (%v vs %v)",
			results[0].UrgencyScore, results[1].UrgencyScore)
	}
	if results[0].Part.ID != "b" {
		t.Errorf("esperado peça 'b' (criticality maior) primeiro, obtido %s", results[0].Part.ID)
	}
}

func TestRankPriorities_DesempatePorAverageDailySales(t *testing.T) {

	a := newTestPart(func(p *Part) {
		p.ID = "a"
		p.Name = "A"
		p.CurrentStock = 10
		p.MinimumStock = 20
		p.AverageDailySales = 9
		p.LeadTimeDays = 0
		p.CriticalityLevel = 2
	})
	b := newTestPart(func(p *Part) {
		p.ID = "b"
		p.Name = "B"
		p.CurrentStock = 10
		p.MinimumStock = 20
		p.AverageDailySales = 1
		p.LeadTimeDays = 0
		p.CriticalityLevel = 2
	})

	results := RankPriorities([]Part{a, b})

	if results[0].UrgencyScore != results[1].UrgencyScore {
		t.Fatalf("pré-condição falhou: scores deveriam ser iguais (%v vs %v)",
			results[0].UrgencyScore, results[1].UrgencyScore)
	}
	if results[0].Part.ID != "a" {
		t.Errorf("esperado peça 'a' (average daily sales maior) primeiro, obtido %s", results[0].Part.ID)
	}
}

func TestRankPriorities_DesempatePorNome(t *testing.T) {

	a := newTestPart(func(p *Part) {
		p.ID = "a"
		p.Name = "Abacate"
		p.CurrentStock = 10
		p.MinimumStock = 20
		p.AverageDailySales = 0
		p.LeadTimeDays = 0
		p.CriticalityLevel = 2
	})
	b := newTestPart(func(p *Part) {
		p.ID = "b"
		p.Name = "Batata"
		p.CurrentStock = 10
		p.MinimumStock = 20
		p.AverageDailySales = 0
		p.LeadTimeDays = 0
		p.CriticalityLevel = 2
	})

	results := RankPriorities([]Part{a, b})

	if results[0].UrgencyScore != results[1].UrgencyScore {
		t.Fatalf("pré-condição falhou: scores deveriam ser iguais (%v vs %v)",
			results[0].UrgencyScore, results[1].UrgencyScore)
	}
	if results[0].Part.ID != "a" {
		t.Errorf("esperado peça 'a' (ordem alfabética) primeiro, obtido %s", results[0].Part.ID)
	}
}

func TestRankPriorities_ListaVazia(t *testing.T) {
	results := RankPriorities([]Part{})
	if len(results) != 0 {
		t.Errorf("esperado lista vazia, obtido %d resultados", len(results))
	}
}
