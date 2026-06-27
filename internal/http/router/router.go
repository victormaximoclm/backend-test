package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"backend-test/internal/http/handler"
)

// New monta o roteador HTTP completo, conectando cada rota ao handler
// correspondente.
func New(partHandler *handler.PartHandler, priorityHandler *handler.PriorityHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/parts", func(r chi.Router) {
		r.Post("/", partHandler.Create)
		r.Get("/", partHandler.List)
		r.Get("/{id}", partHandler.Get)
		r.Put("/{id}", partHandler.Update)
		r.Delete("/{id}", partHandler.Delete)
	})

	r.Route("/restock", func(r chi.Router) {
		r.Get("/priorities", priorityHandler.GetPriorities)
	})

	return r
}
