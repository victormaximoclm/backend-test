package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-test/internal/http/handler"
	"backend-test/internal/http/router"
	"backend-test/internal/repository"
	"backend-test/internal/service"
)

// newTestServer monta a aplicação completa (repositório real + services
// reais + handlers + router), exatamente como o main.go faz, para validar
// o fluxo HTTP de ponta a ponta.
func newTestServer() http.Handler {
	repo := repository.NewInMemoryPartRepository()
	partSvc := service.NewPartService(repo)
	prioritySvc := service.NewPriorityService(repo)

	partHandler := handler.NewPartHandler(partSvc)
	priorityHandler := handler.NewPriorityHandler(prioritySvc)

	return router.New(partHandler, priorityHandler)
}

func TestIntegration_HealthCheck(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d; esperado 200", rec.Code)
	}
}

func TestIntegration_CreateAndGetPart(t *testing.T) {
	srv := newTestServer()

	payload := handler.PartRequest{
		Name:              "Filtro de Óleo X",
		Category:          "engine",
		CurrentStock:      15,
		MinimumStock:      20,
		AverageDailySales: 4,
		LeadTimeDays:      5,
		UnitCost:          18.50,
		CriticalityLevel:  3,
	}
	body, _ := json.Marshal(payload)

	createReq := httptest.NewRequest(http.MethodPost, "/parts/", bytes.NewReader(body))
	createRec := httptest.NewRecorder()
	srv.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("status = %d; esperado 201. body: %s", createRec.Code, createRec.Body.String())
	}

	var created handler.PartResponse
	json.Unmarshal(createRec.Body.Bytes(), &created)
	if created.ID == "" {
		t.Fatal("esperado ID gerado, obtido vazio")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/parts/"+created.ID, nil)
	getRec := httptest.NewRecorder()
	srv.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("status = %d; esperado 200", getRec.Code)
	}
}

func TestIntegration_RestockPriorities(t *testing.T) {
	srv := newTestServer()

	// Peça que precisa de reposição urgente (exemplo do desafio).
	payload := handler.PartRequest{
		Name:              "Filtro de Óleo X",
		Category:          "engine",
		CurrentStock:      15,
		MinimumStock:      20,
		AverageDailySales: 4,
		LeadTimeDays:      5,
		UnitCost:          18.50,
		CriticalityLevel:  3,
	}
	body, _ := json.Marshal(payload)
	createReq := httptest.NewRequest(http.MethodPost, "/parts/", bytes.NewReader(body))
	createRec := httptest.NewRecorder()
	srv.ServeHTTP(createRec, createReq)

	// Peça com estoque confortável, não deve aparecer no resultado.
	calmPayload := handler.PartRequest{
		Name:              "Peça Tranquila",
		Category:          "engine",
		CurrentStock:      1000,
		MinimumStock:      10,
		AverageDailySales: 1,
		LeadTimeDays:      1,
		UnitCost:          5.0,
		CriticalityLevel:  1,
	}
	calmBody, _ := json.Marshal(calmPayload)
	calmReq := httptest.NewRequest(http.MethodPost, "/parts/", bytes.NewReader(calmBody))
	calmRec := httptest.NewRecorder()
	srv.ServeHTTP(calmRec, calmReq)

	// Chama o endpoint de priorização.
	prioReq := httptest.NewRequest(http.MethodGet, "/restock/priorities", nil)
	prioRec := httptest.NewRecorder()
	srv.ServeHTTP(prioRec, prioReq)

	if prioRec.Code != http.StatusOK {
		t.Fatalf("status = %d; esperado 200. body: %s", prioRec.Code, prioRec.Body.String())
	}

	var resp handler.PrioritiesResponse
	json.Unmarshal(prioRec.Body.Bytes(), &resp)

	if len(resp.Priorities) != 1 {
		t.Fatalf("esperado 1 peça na priorização, obtido %d", len(resp.Priorities))
	}
	if resp.Priorities[0].Name != "Filtro de Óleo X" {
		t.Errorf("nome incorreto: %s", resp.Priorities[0].Name)
	}
	if resp.Priorities[0].UrgencyScore != 75 {
		t.Errorf("urgencyScore = %v; esperado 75", resp.Priorities[0].UrgencyScore)
	}
}

func TestIntegration_CreatePart_ValidacaoFalha(t *testing.T) {
	srv := newTestServer()

	payload := handler.PartRequest{
		Name:             "Peça Inválida",
		Category:         "engine",
		CriticalityLevel: 10, // fora do range 1-5
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/parts/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d; esperado 400. body: %s", rec.Code, rec.Body.String())
	}
}

func TestIntegration_GetPart_NaoEncontrada(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/parts/id-que-nao-existe", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d; esperado 404", rec.Code)
	}
}
