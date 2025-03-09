// repository/generic_repo.go
package repository

type Repository[T any] interface {
	GetByID(id int64) (*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id int64) error
}
