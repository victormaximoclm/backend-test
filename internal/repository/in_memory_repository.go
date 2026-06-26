package repository

import (
	"backend-test/internal/domain"
	"sync"
	"time"
)

// Estrutura do repositório em memória.
// Com controle de acesso concorrente e um map que armazena as peças pelo ID.
type InMemoryPartRepository struct {
	mu    sync.RWMutex
	parts map[string]domain.Part
}

// NewInMemoryPartRepository cria um repositório em memória vazio
func NewInMemoryPartRepository() *InMemoryPartRepository {
	return &InMemoryPartRepository{
		parts: make(map[string]domain.Part),
	}
}

// Create adiciona uma peça ao repositório em memória.
func (r *InMemoryPartRepository) Create(part domain.Part) (domain.Part, error) {
	r.mu.Lock()         //trava para acessos simultaneos
	defer r.mu.Unlock() //destrava quando terminar

	now := time.Now().UTC()
	part.CreatedAt = now
	part.UpdatedAt = now

	r.parts[part.ID] = part
	return part, nil
}

// GetByID busca uma peça pelo ID no repositório.
// Retorna a peça encontrada ou erro caso não exista.
func (r *InMemoryPartRepository) GetByID(id string) (domain.Part, error) {
	r.mu.RLock() // bloqueia escrita enquanto realiza leitura do map
	defer r.mu.RUnlock()

	part, exists := r.parts[id]
	if !exists {
		return domain.Part{}, ErrPartNotFound
	}
	return part, nil
}

// List percorre o map interno e retorna todas as peças em um slice.
func (r *InMemoryPartRepository) List() ([]domain.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Part, 0, len(r.parts)) //cria slice do tamanho de quantidade de pecas do repositorio
	for _, p := range r.parts {                    //percorre ate o final do repositorio
		result = append(result, p) //adiciona peca ao slice
	}
	return result, nil //retorna o slice preenchido
}

// ListByCategory funciona de forma semlhante mas antes verifica se a categoria da peca equivale a solicitada
func (r *InMemoryPartRepository) ListByCategory(category string) ([]domain.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Part, 0)
	for _, p := range r.parts {
		if p.Category == category {
			result = append(result, p)
		}
	}
	return result, nil
}

// Update Atualiza informações da peça
func (r *InMemoryPartRepository) Update(part domain.Part) (domain.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.parts[part.ID] //verifica se a peça existe através do ID da struct
	if !exists {
		return domain.Part{}, ErrPartNotFound
	}

	part.CreatedAt = existing.CreatedAt
	part.UpdatedAt = time.Now().UTC()

	r.parts[part.ID] = part //recebe nova struct da peça
	return part, nil
}

// Remove a entrada correspondente ao ID do map.
func (r *InMemoryPartRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.parts[id]; !exists {
		return ErrPartNotFound
	}

	delete(r.parts, id) //deleta a chave
	return nil
}
