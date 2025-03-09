package services

import (
	"reflect"
	"sync"

	"github.com/PunyGoood/DCS/system/repository"
)

type ModelService[T any] struct {
	repo repository.Repository[T]
}

var (
	mu         sync.Mutex
	serviceMap = make(map[string]interface{})
)

func NewModelService[T any](repo repository.Repository[T]) *ModelService[T] {
	typeName := reflect.TypeOf(new(T)).Elem().String()

	mu.Lock()
	defer mu.Unlock()

	if instance, exists := serviceMap[typeName]; exists {
		return instance.(*ModelService[T])
	}

	service := &ModelService[T]{repo: repo}
	serviceMap[typeName] = service
	return service
}

func (svc *ModelService[T]) GetById(id int64) (*T, error) {
	return svc.repo.GetByID(id)
}
