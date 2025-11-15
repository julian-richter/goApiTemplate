package repos

import (
	"fmt"

	"github.com/julian-richter/ApiTemplate/internal/models"
)

// Repository provides generic CRUD operations for any Entity type.
// It uses Go generics to maintain strong typing without reflection.
type Repository[T models.Entity] struct {
	data map[int]T
}

func NewRepository[T models.Entity]() *Repository[T] {
	return &Repository[T]{data: make(map[int]T)}
}

// Save stores or updates an entity by its ID.
func (r *Repository[T]) Save(entity T) {
	r.data[entity.GetID()] = entity
}

// GetByID retrieves an entity by its ID.
func (r *Repository[T]) GetByID(id int) (T, error) {
	entity, exists := r.data[id]
	if !exists {
		var zero T
		return zero, fmt.Errorf("entity with id %d not found", id)
	}
	return entity, nil
}

// Delete removes an entity from the repository.
func (r *Repository[T]) Delete(id int) {
	delete(r.data, id)
}

// All returns all entities in the repository.
func (r *Repository[T]) All() []T {
	result := make([]T, 0, len(r.data))
	for _, e := range r.data {
		result = append(result, e)
	}
	return result
}
