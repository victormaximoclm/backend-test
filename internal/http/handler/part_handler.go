package handler

import (
	"backend-test/internal/domain"
	"backend-test/internal/repository"
	"backend-test/internal/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// statusForError mapeia erros de domínio/repositório para status HTTP apropriados
func statusForError(err error) int {
	switch {
	case errors.Is(err, repository.ErrPartNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrEmptyName),
		errors.Is(err, domain.ErrEmptyCategory),
		errors.Is(err, domain.ErrNegativeMinimumStock),
		errors.Is(err, domain.ErrNegativeDailySales),
		errors.Is(err, domain.ErrNegativeLeadTime),
		errors.Is(err, domain.ErrNegativeUnitCost),
		errors.Is(err, domain.ErrInvalidCriticality):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// PartHandler traduz requisições HTTP em chamadas ao PartService
type PartHandler struct {
	service *service.PartService
}

func NewPartHandler(s *service.PartService) *PartHandler {
	return &PartHandler{service: s}
}

func (h *PartHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req PartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido: "+err.Error())
		return
	}

	part, err := h.service.CreatePart(req.toCreateInput())
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, newPartResponse(part))
}

func (h *PartHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	part, err := h.service.GetPart(id)
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, newPartResponse(part))
}

func (h *PartHandler) List(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var parts []domain.Part
	var err error

	if category != "" {
		parts, err = h.service.ListPartsByCategory(category)
	} else {
		parts, err = h.service.ListParts()
	}

	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, newPartListResponse(parts))
}

func (h *PartHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req PartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido: "+err.Error())
		return
	}

	part, err := h.service.UpdatePart(id, req.toUpdateInput())
	if err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, newPartResponse(part))
}

func (h *PartHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeletePart(id); err != nil {
		writeError(w, statusForError(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
