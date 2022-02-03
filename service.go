package di

import "reflect"

type Service struct {
	definition reflect.Type
	impl       reflect.Type

	value *reflect.Value

	dependencies []*Service

	setup    func(c *Container) (any, error)
	teardown func(c *Container) error
}

func (t *Service) Teardown(c *Container) error {
	if t.teardown != nil {
		return t.teardown(c)
	}

	if t.value != nil {
		if td, ok := t.value.Interface().(Teardown); ok {
			return td.Teardown()
		}
	}

	return nil
}

type Teardown interface {
	Teardown() error
}
