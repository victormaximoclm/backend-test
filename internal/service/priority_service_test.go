package service

import "testing"

func TestPriorityService_GetPriorities_OrdenaCorretamente(t *testing.T) {
	repo := newFakePartRepository()
	prioritySvc := NewPriorityService(repo)
	partSvc := NewPartService(repo)

	// Peça com alta urgência (estoque negativo).
	partSvc.CreatePart(CreatePartInput{
		Name:              "Pastilha de Freio Y",
		Category:          "brakes",
		CurrentStock:      -2,
		MinimumStock:      10,
		AverageDailySales: 2,
		LeadTimeDays:      3,
		UnitCost:          45.0,
		CriticalityLevel:  5,
	})

	// Peça sem necessidade de reposição.
	partSvc.CreatePart(CreatePartInput{
		Name:              "Peça Tranquila",
		Category:          "engine",
		CurrentStock:      500,
		MinimumStock:      10,
		AverageDailySales: 1,
		LeadTimeDays:      2,
		UnitCost:          5.0,
		CriticalityLevel:  1,
	})

	// Peça que precisa reposição, urgência intermediária.
	partSvc.CreatePart(CreatePartInput{
		Name:              "Filtro de Óleo X",
		Category:          "engine",
		CurrentStock:      15,
		MinimumStock:      20,
		AverageDailySales: 4,
		LeadTimeDays:      5,
		UnitCost:          18.50,
		CriticalityLevel:  3,
	})

	results, err := prioritySvc.GetPriorities()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("esperado 2 peças precisando de reposição, obtido %d", len(results))
	}
	if results[0].Part.Name != "Pastilha de Freio Y" {
		t.Errorf("esperado 'Pastilha de Freio Y' como mais urgente, obtido %q", results[0].Part.Name)
	}
	if results[0].UrgencyScore < results[1].UrgencyScore {
		t.Errorf("resultados fora de ordem: %v < %v", results[0].UrgencyScore, results[1].UrgencyScore)
	}
}

func TestPriorityService_GetPriorities_SemPecas(t *testing.T) {
	repo := newFakePartRepository()
	svc := NewPriorityService(repo)

	results, err := svc.GetPriorities()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("esperado lista vazia, obtido %d resultados", len(results))
	}
}
