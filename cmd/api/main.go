package main

import (
	"log"
	"net/http"
	"os"

	"backend-test/internal/http/handler"
	"backend-test/internal/http/router"
	"backend-test/internal/repository"
	"backend-test/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Camada de persistência. Trocar de banco no futuro significa
	// substituir esta linha por outra implementação de
	// repository.PartRepository — nenhuma linha abaixo precisa mudar.
	partRepo := repository.NewInMemoryPartRepository()

	// Camada de serviço.
	partService := service.NewPartService(partRepo)
	priorityService := service.NewPriorityService(partRepo)

	// Camada HTTP.
	partHandler := handler.NewPartHandler(partService)
	priorityHandler := handler.NewPriorityHandler(priorityService)

	r := router.New(partHandler, priorityHandler)

	addr := ":" + port
	log.Printf("restock-priority-service ouvindo em %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
