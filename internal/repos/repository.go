package repos

import (
	"fmt"
	"sync"

	"github.com/julian-richter/ApiTemplate/internal/models"
)

type Repository[T models.Entity] struct {
	mu   sync.RWMutex
	data map[int]T
}

func NewRepository[T models.Entity]() *Repository[T] {
	return &Repository[T]{
		data: make(map[int]T),
	}
}

// ensureMap lazily initializes r.data if the repository was zero-valued.
func (r *Repository[T]) ensureMap() {
	if r.data == nil {
		r.data = make(map[int]T)
	}
}

func (r *Repository[T]) Save(entity T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ensureMap()
	r.data[int(entity.GetID())] = entity
}

func (r *Repository[T]) GetByID(id int) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var zero T
	entity, ok := r.data[id]
	if !ok {
		return zero, fmt.Errorf("entity with id %d not found", id)
	}
	return entity, nil
}

func (r *Repository[T]) Delete(id int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, id)
}

func (r *Repository[T]) All() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]T, 0, len(r.data))
	for _, e := range r.data {
		result = append(result, e)
	}
	return result
}
