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
	if !tImpl.Implements(tIf) && !reflect.PointerTo(tImpl).Implements(tIf) {
		err = ErrDoesNotImplInterface
		return
	}
	key := getInterfaceKey(tIf)
	c.Put(key, &singletonService{
		ImplType: tImpl,
	})
	return
}

func MustRegister[TSvc, TImpl any](c Container) {
	must(Register[TSvc, TImpl](c))
}

// RegisterTransient adds a transient service of the type TSvc
// with an implementation of the type TImpl to the provided container.
func RegisterTransient[TSvc, TImpl any](c Container) (err error) {
	tImpl := getType[TImpl]()
	tIf := getType[TSvc]()
	if tIf.Kind() != reflect.Interface {
		err = ErrNoInterface
		return
	}
	if !tImpl.Implements(tIf) && !reflect.PointerTo(tImpl).Implements(tIf) {
		err = ErrDoesNotImplInterface
		return
	}
	key := getInterfaceKey(tIf)
	c.Put(key, &transientService{
		implType: tImpl,
	})
	return
}

// MustRegisterTransient adds a transient service of the type TSvc
// with an implementation of the type TImpl to the provided container.
//
// It panics when a registration is not possible.
func MustRegisterTransient[TSvc, TImpl any](c Container) {
	must(RegisterTransient[TSvc, TImpl](c))
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

func MustGet[T any](c Container) (s T) {
	s, err := Get[T](c)
	must(err)
	return
}
