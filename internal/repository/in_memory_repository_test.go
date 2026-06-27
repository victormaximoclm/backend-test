package repository

import (
	"fmt"
	"sync"
	"testing"

	"backend-test/internal/domain"
)

// Cria struct de peça teste
func testPart(id string) domain.Part {
	return domain.Part{
		ID:                id,
		Name:              "Peça Teste",
		Category:          "engine",
		CurrentStock:      15,
		MinimumStock:      20,
		AverageDailySales: 4,
		LeadTimeDays:      5,
		UnitCost:          18.50,
		CriticalityLevel:  3,
	}
}

// Cria repositório vazio, adiciona peça teste ao repositorio e busca peça
func TestInMemoryPartRepository_CreateAndGetByID(t *testing.T) {
	repo := NewInMemoryPartRepository()

	created, err := repo.Create(testPart("p1"))
	if err != nil {
		t.Fatalf("erro inesperado ao criar: %v", err) //para o teste se não conseguir criar
	}
	if created.CreatedAt.IsZero() {
		t.Error("CreatedAt deveria ser preenchido pelo repositório") //retorna erro se createdAt não vir preenchido
	}

	found, err := repo.GetByID("p1")
	if err != nil {
		t.Fatalf("erro inesperado ao buscar: %v", err) //para o teste e retorna erro por não achar
	}
	if found.ID != "p1" {
		t.Errorf("ID = %s; esperado p1", found.ID) //avisa erro caso ache uma peça mas não seja a que era esperada
	}
}

// testa o caso de não achar a peça dentro do repositorio
func TestInMemoryPartRepository_GetByID_NaoEncontrado(t *testing.T) {
	repo := NewInMemoryPartRepository()

	_, err := repo.GetByID("inexistente") //ignora a peça e armazena o erro
	if err != ErrPartNotFound {
		t.Errorf("esperado ErrPartNotFound, obtido: %v", err) //retorna erro caso o erro não seja de peça não encontrada
	}
}

// testa a listagem do repositorio
func TestInMemoryPartRepository_List(t *testing.T) {
	repo := NewInMemoryPartRepository()
	repo.Create(testPart("p1"))
	repo.Create(testPart("p2")) //cria duas peças e armazena no repositorio

	all, err := repo.List()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("esperado 2 peças, obtido %d", len(all)) //retorna erro se o tamanho não é o esperado
	}
}

// testa a listagem por categoria
func TestInMemoryPartRepository_ListByCategory(t *testing.T) {
	repo := NewInMemoryPartRepository()

	engine := testPart("p1")   //pega estrutura teste
	engine.Category = "engine" //muda categoria
	repo.Create(engine)

	brakes := testPart("p2")
	brakes.Category = "brakes"
	repo.Create(brakes)

	result, err := repo.ListByCategory("brakes") //lista peças da categoria ''brakes''
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("esperado 1 peça na categoria 'brakes', obtido %d", len(result))
	}
	if result[0].ID != "p2" {
		t.Errorf("peça retornada incorreta: %s", result[0].ID)
	}
}

// Testa a atualização de peça
func TestInMemoryPartRepository_Update(t *testing.T) {
	repo := NewInMemoryPartRepository()
	created, _ := repo.Create(testPart("p1"))

	toUpdate := created
	toUpdate.CurrentStock = 999 //atualiza a quantidade de currentStock

	updated, err := repo.Update(toUpdate) //lança o comando de atualização
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if updated.CurrentStock != 999 {
		t.Errorf("CurrentStock = %d; esperado 999", updated.CurrentStock)
	}
	// CreatedAt deve ser preservado do registro original, não sobrescrito.
	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("CreatedAt não deveria mudar em um Update")
	}
}

// testa o update caso não encontre a peça
func TestInMemoryPartRepository_Update_NaoEncontrado(t *testing.T) {
	repo := NewInMemoryPartRepository()

	_, err := repo.Update(testPart("inexistente"))
	if err != ErrPartNotFound {
		t.Errorf("esperado ErrPartNotFound, obtido: %v", err)
	}
}

// testa o delete da peça
func TestInMemoryPartRepository_Delete(t *testing.T) {
	repo := NewInMemoryPartRepository()
	repo.Create(testPart("p1"))

	if err := repo.Delete("p1"); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	_, err := repo.GetByID("p1")
	if err != ErrPartNotFound {
		t.Errorf("peça deveria ter sido removida")
	}
}

// testa caso a peça a ser deletada não exista
func TestInMemoryPartRepository_Delete_NaoEncontrado(t *testing.T) {
	repo := NewInMemoryPartRepository()

	err := repo.Delete("inexistente")
	if err != ErrPartNotFound {
		t.Errorf("esperado ErrPartNotFound, obtido: %v", err)
	}
}

// Teste de concorrência
//
// Valida que o repositório suporta acesso simultâneo sem corromper dados
// para suportar centenas ou milhares de peças
//
//	go test ./internal/repository/... -race -v
func TestInMemoryPartRepository_ConcorrenciaCriacaoEListagem(t *testing.T) {
	repo := NewInMemoryPartRepository()

	const totalGoroutines = 200
	var wg sync.WaitGroup

	for i := 0; i < totalGoroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("id-%d", n)
			_, err := repo.Create(testPart(id))
			if err != nil {
				t.Errorf("erro inesperado ao criar peça %s: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	all, err := repo.List()
	if err != nil {
		t.Fatalf("erro inesperado ao listar: %v", err)
	}
	if len(all) != totalGoroutines {
		t.Errorf("esperado %d peças, obtido %d", totalGoroutines, len(all))
	}
}

// TestInMemoryPartRepository_ConcorrenciaLeituraEEscrita mistura leituras
// (List) e escritas (Create) simultâneas, para validar que o RWMutex
// protege corretamente os dois tipos de acesso ao mesmo tempo.
func TestInMemoryPartRepository_ConcorrenciaLeituraEEscrita(t *testing.T) {
	repo := NewInMemoryPartRepository()

	for i := 0; i < 50; i++ {
		repo.Create(testPart(fmt.Sprintf("seed-%d", i)))
	}

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = repo.List()
		}()
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_, _ = repo.Create(testPart(fmt.Sprintf("new-%d", n)))
		}(i)
	}

	wg.Wait()

	all, err := repo.List()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(all) != 100 {
		t.Errorf("esperado 100 peças, obtido %d", len(all))
	}
}
