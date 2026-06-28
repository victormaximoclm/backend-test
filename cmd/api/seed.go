package main

import (
	"log"

	"backend-test/internal/domain"
	"backend-test/internal/repository"
)

// seedData popula o repositório com peças de exemplo, para facilitar
// testes manuais imediatos após iniciar o servidor (sem precisar criar
// peças via POST antes de poder testar GET/PUT/DELETE e a priorização).
//
// Usa o repositório diretamente (não o service), porque queremos
// controlar o ID das peças de exemplo
func seedData(repo *repository.InMemoryPartRepository) {
	seedParts := []domain.Part{
		{
			ID:                "11111111-1111-1111-1111-111111111111",
			Name:              "Filtro de Óleo X",
			Category:          "engine",
			CurrentStock:      15,
			MinimumStock:      20,
			AverageDailySales: 4,
			LeadTimeDays:      5,
			UnitCost:          18.50,
			CriticalityLevel:  3,
		},
		{
			ID:                "22222222-2222-2222-2222-222222222222",
			Name:              "Pastilha de Freio Y",
			Category:          "brakes",
			CurrentStock:      -2,
			MinimumStock:      10,
			AverageDailySales: 2,
			LeadTimeDays:      3,
			UnitCost:          45.00,
			CriticalityLevel:  5,
		},
		{
			ID:                "33333333-3333-3333-3333-333333333333",
			Name:              "Vela de Ignição Z",
			Category:          "engine",
			CurrentStock:      500,
			MinimumStock:      10,
			AverageDailySales: 1,
			LeadTimeDays:      2,
			UnitCost:          12.00,
			CriticalityLevel:  1,
		},
		{
			ID:                "44444444-4444-4444-4444-444444444444",
			Name:              "Amortecedor Traseiro",
			Category:          "suspension",
			CurrentStock:      8,
			MinimumStock:      15,
			AverageDailySales: 1.5,
			LeadTimeDays:      10,
			UnitCost:          120.00,
			CriticalityLevel:  4,
		},
		{
			ID:                "55555555-5555-5555-5555-555555555555",
			Name:              "Correia Dentada",
			Category:          "engine",
			CurrentStock:      12,
			MinimumStock:      20,
			AverageDailySales: 0,
			LeadTimeDays:      7,
			UnitCost:          65.00,
			CriticalityLevel:  4,
		},
	}

	for _, p := range seedParts {
		if _, err := repo.Create(p); err != nil {
			log.Printf("erro ao criar peça de seed %s: %v", p.Name, err)
		}
	}

	log.Printf("seed concluido: %d pecas de exemplo criadas", len(seedParts))
}
