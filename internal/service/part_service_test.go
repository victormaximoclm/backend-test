package service

import (
	"errors"
	"testing"

	"backend-test/internal/domain"
	"backend-test/internal/repository"
)

func validCreateInput() CreatePartInput {
	return CreatePartInput{
		Name:              "Filtro de Óleo X",
		Category:          "engine",
		CurrentStock:      15,
		MinimumStock:      20,
		AverageDailySales: 4,
		LeadTimeDays:      5,
		UnitCost:          18.50,
		CriticalityLevel:  3,
	}
}

func TestPartService_CreatePart_Sucesso(t *testing.T) {
	repo := newFakePartRepository()
	svc := NewPartService(repo)

	part, err := svc.CreatePart(validCreateInput())
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if part.ID == "" {
		t.Error("esperado ID gerado automaticamente, obtido vazio")
	}
	if part.Name != "Filtro de Óleo X" {
		t.Errorf("name = %q; esperado 'Filtro de Óleo X'", part.Name)
	}
}

func TestPartService_CreatePart_ValidacaoFalha(t *testing.T) {
	repo := newFakePartRepository()
	svc := NewPartService(repo)

	input := validCreateInput()
	input.CriticalityLevel = 99 // inválido, fora do range 1-5

	_, err := svc.CreatePart(input)
	if !errors.Is(err, domain.ErrInvalidCriticality) {
		t.Errorf("esperado ErrInvalidCriticality, obtido: %v", err)
	}
}

func TestPartService_GetPart_NaoEncontrada(t *testing.T) {
	repo := newFakePartRepository()
	svc := NewPartService(repo)

	_, err := svc.GetPart("id-inexistente")
	if !errors.Is(err, repository.ErrPartNotFound) {
		t.Errorf("esperado ErrPartNotFound, obtido: %v", err)
	}
}

func TestPartService_UpdatePart_Sucesso(t *testing.T) {
	repo := newFakePartRepository()
	svc := NewPartService(repo)

	created, _ := svc.CreatePart(validCreateInput())

	updated, err := svc.UpdatePart(created.ID, UpdatePartInput{
		Name:              "Filtro de Óleo X Atualizado",
		Category:          "engine",
		CurrentStock:      30,
		MinimumStock:      20,
		AverageDailySales: 4,
		LeadTimeDays:      5,
		UnitCost:          19.90,
		CriticalityLevel:  3,
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if updated.Name != "Filtro de Óleo X Atualizado" {
		t.Errorf("name não foi atualizado, obtido: %q", updated.Name)
	}
	if updated.ID != created.ID {
		t.Errorf("ID não deveria mudar em um update")
	}
}

func TestPartService_UpdatePart_NaoEncontrada(t *testing.T) {
	repo := newFakePartRepository()
	svc := NewPartService(repo)

	_, err := svc.UpdatePart("id-inexistente", UpdatePartInput{
		Name:             "X",
		Category:         "engine",
		CriticalityLevel: 1,
	})
	if !errors.Is(err, repository.ErrPartNotFound) {
		t.Errorf("esperado ErrPartNotFound, obtido: %v", err)
	}
}

func TestPartService_DeletePart_Sucesso(t *testing.T) {
	repo := newFakePartRepository()
	svc := NewPartService(repo)

	created, _ := svc.CreatePart(validCreateInput())

	if err := svc.DeletePart(created.ID); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	_, err := svc.GetPart(created.ID)
	if !errors.Is(err, repository.ErrPartNotFound) {
		t.Errorf("peça deveria ter sido removida, mas ainda foi encontrada")
	}
}

func TestPartService_ListPartsByCategory(t *testing.T) {
	repo := newFakePartRepository()
	svc := NewPartService(repo)

	engineInput := validCreateInput()
	engineInput.Category = "engine"
	svc.CreatePart(engineInput)

	brakeInput := validCreateInput()
	brakeInput.Name = "Pastilha de Freio Y"
	brakeInput.Category = "brakes"
	svc.CreatePart(brakeInput)

	results, err := svc.ListPartsByCategory("brakes")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("esperado 1 resultado na categoria 'brakes', obtido %d", len(results))
	}
	if results[0].Category != "brakes" {
		t.Errorf("categoria incorreta retornada: %s", results[0].Category)
	}
}
