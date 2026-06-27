package handler

import (
	"encoding/json"
	"net/http"
)

// writeJSON serializa qualquer payload como JSON com o status informado.
// Centralizado aqui para garantir Content-Type consistente em toda a API.
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError padroniza o formato de erro retornado ao cliente.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
