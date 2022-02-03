package di

import (
	"reflect"
)

type Service struct {
	ImplType reflect.Type
	IsBuilt  bool
	Instance reflect.Value
}

func (s *Service) Build(c Container) (instance reflect.Value) {
	if s.IsBuilt {
		instance = s.Instance
		return
	}
	instance = reflect.New(s.ImplType)
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
		fInstance := svc.Build(c)
		tF.Set(fInstance)
	}
	s.Instance = instance
	s.IsBuilt = true
	return
}
