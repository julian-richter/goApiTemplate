package models

// Entity defines the basic behavior for all storable models.
// It allows repositories to manage entities generically while maintaining type safety.
type Entity interface {
	GetID() int
	SetID(id int)
}
