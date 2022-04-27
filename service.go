package di

import (
	"reflect"
)

// Service describes what type of object should get constructed
// and how should it get returned.
type Service interface {
	Build(c Container) (instance reflect.Value)
}

// singletonService describes a dependency that should only be constructed once
// and reused for the lifetime of the application.
type singletonService struct {
	ImplType reflect.Type
	IsBuilt  bool
	Instance reflect.Value
}

func (s *singletonService) Build(c Container) (instance reflect.Value) {
	if s.IsBuilt {
		instance = s.Instance
		return
	}
	instance = reflect.New(s.ImplType)
	s.Instance = instance
	s.IsBuilt = true
	injectFields(instance, c)
	return
}

// transientService describes a dependency that should get created anew
// each time it is requested from the container.
type transientService struct {
	implType reflect.Type
}

func (s *transientService) Build(c Container) (instance reflect.Value) {
	instance = reflect.New(s.implType)
	injectFields(instance, c)
	return
}

// injectFields attempts to inject dependencies into a given value,
// using a provided DI container.
func injectFields(instance reflect.Value, c Container) {
	elem := instance.Elem()
	for i := 0; i < elem.NumField(); i++ {
		tF := elem.Field(i)
		if tF.Kind() != reflect.Interface {
			continue
		}
		key := getInterfaceKey(tF.Type())
		svc, ok := c.Get(key)
		if !ok {
			continue
		}
		if tF.CanSet() && tF.IsNil() {
			fInstance := svc.Build(c)
			tF.Set(fInstance)
		}
	}
}
