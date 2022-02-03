package di

import (
	"reflect"
)

func Register[TSvc, TImpl any](c Container) (err error) {
	tImpl := getType[TImpl]()
	tIf := getType[TSvc]()
	if tIf.Kind() != reflect.Interface {
		err = ErrNoInterface
		return
	}
	if !tImpl.Implements(tIf) {
		err = ErrDoesNotImplInterface
		return
	}
	key := getInterfaceKey(tIf)
	c.Put(key, &Service{
		ImplType: tImpl,
	})
	return
}

func Get[T any](c Container) (s T, err error) {
	tIf := getType[T]()
	if tIf.Kind() != reflect.Interface {
		err = ErrNoInterface
		return
	}
	key := getInterfaceKey(tIf)
	sb, ok := c.Get(key)
	if !ok {
		err = ErrNotRegistered
		return
	}
	v := sb.Build(c).Interface()
	s, ok = v.(T)
	if !ok {
		err = ErrInvalidImplementation
	}
	return
}
