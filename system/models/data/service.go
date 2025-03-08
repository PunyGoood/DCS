package data

import (
	"reflect"
	"sync"
)

var (
	mu          sync.Mutex
	onceMap     = make(map[string]*sync.Once)
	instanceMap = make(map[string]interface{})
)
type ModelServiceV2[T any] struct {
	
}
func NewModelService[T any]() *ModelService[T] {
	typeName := reflect.TypeOf(new(T)).Elem().String() // 更可靠的类型名获取方式

	mu.Lock()

	once, exists := onceMap[typeName] // onceMap[string]*sync.Once
	if !exists {
		once = new(sync.Once)
		onceMap[typeName] = once
		instanceMap[typeName] = nil
	}
	mu.Unlock()

	var instance *ModelService[T]
	once.Do(func() {
		collectionName := getCollectionName[T]()
		
		// 建立数据库连接   待定
		collection := 

		instance = 

		mu.Lock()
		instanceMap[typeName] = instance
		mu.Unlock()
	})

	mu.Lock()
	defer mu.Unlock()
	return instanceMap[typeName].(*ModelService[T])
}

func (svc *ModelService[T]) GetById(id int64) (model *T, err error) {
	var result T
	err = svc.col.FindId(id).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
