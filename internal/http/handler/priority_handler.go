package handler

import (
	"net/http"

	"backend-test/internal/service"
)

// PriorityHandler expõe o endpoint GET /restock/priorities.
type PriorityHandler struct {
	service *service.PriorityService
}

func NewPriorityHandler(s *service.PriorityService) *PriorityHandler {
	return &PriorityHandler{service: s}
}

func (h *PriorityHandler) GetPriorities(w http.ResponseWriter, r *http.Request) {
	results, err := h.service.GetPriorities()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, newPrioritiesResponse(results))
}
