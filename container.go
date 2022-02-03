package di

import (
	"errors"
	"reflect"
	"sync"
)

// Container is used to register and retrieve
// dependency instances and control their
// lifecycle.
type Container struct {
	mx       sync.Mutex
	services map[string]*Service
}

// RegisterOptions defines options for service
// registration.
type RegisterOptions[TImpl any] struct {
	Setup    func(c *Container) (*TImpl, error)
	Teardown func(c *Container) error
	Instance *TImpl
}

// NewContainer returns a new sinatcne
// of Container.
func NewContainer() *Container {
	return &Container{
		services: make(map[string]*Service),
	}
}

// Teardown executes the teardown functions
// or methods, if available, of all initiated
// services in reversed order of the
// dependency tree.
//
// Currently, teardown only works on
// singleton instances.
func (c *Container) Teardown() error {
	var (
		errs     []error
		services []*Service
		queue    []*Service
	)

	c.mx.Lock()
	defer c.mx.Unlock()

	for _, s := range c.services {
		queue = append(queue, s)
	}

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]

		services = append(services, next)
		queue = append(queue, next.dependencies...)
	}

	serviceSet := hashSet[*Service]{}
	for i := len(services) - 1; i >= 0; i-- {
		if s := services[i]; serviceSet.Set(s) {
			errs = append(errs, s.Teardown(c))
		}
	}

	return errors.Join(errs...)
}

// Register takes an instance of Container and registers
// the given TImpl type (as struct) for the specified
// TService type (interface).
//
// When no options are passed, an empty instance of
// TImpl will be created. Otherwise, the given Instance
// or Setup result will be taken (in that order).
func Register[TService, TImpl any](c *Container, opt ...RegisterOptions[TImpl]) error {
	tSvc := getType[TService]()
	tImpl := getType[TImpl]()

	if tSvc.Kind() != reflect.Interface {
		return ErrSvcNoInterface
	}
	if tImpl.Kind() != reflect.Struct {
		return ErrImplNoStruct
	}
	if !tImpl.Implements(tSvc) && !reflect.PointerTo(tImpl).Implements(tSvc) {
		return ErrImplDoesNotImplementSvc
	}

	svc := &Service{
		definition: tSvc,
		impl:       tImpl,
	}

	if len(opt) > 0 {
		if instance := opt[0].Instance; instance != nil {
			svc.setup = func(c *Container) (any, error) {
				return instance, nil
			}
		} else if setup := opt[0].Setup; setup != nil {
			svc.setup = func(c *Container) (any, error) {
				return setup(c)
			}
		}
		if teardown := opt[0].Teardown; teardown != nil {
			svc.teardown = func(c *Container) error {
				return teardown(c)
			}
		}
	}

	key := getInterfaceKey(tSvc)

	c.mx.Lock()
	defer c.mx.Unlock()
	c.services[key] = svc

	return nil
}

// Get retireves or constructs the registered implementation
// for the given TService interface.
//
// When Singleton is used as strategy and an instance has
// already been created for this service, the instance is
// reused. Otherwise, when Transistent is chosen, a new
// instance will be created. The default strategy is
// Singleton when not further specified.
func Get[TService any](c *Container, strategy ...Strategy) (s TService, err error) {
	tSvc := getType[TService]()

	strat := Singleton
	if len(strategy) > 0 {
		strat = strategy[0]
	}

	if tSvc.Kind() != reflect.Interface {
		return s, ErrSvcNoInterface
	}

	key := getInterfaceKey(tSvc)

	c.mx.Lock()
	defer c.mx.Unlock()

	val, err := c.build(key, strat, nil)
	if err != nil {
		return s, err
	}

	s, ok := val.Interface().(TService)
	if !ok {
		return s, ErrStoredTypeCanNotBeCasted
	}

	return s, nil
}

func (c *Container) build(key string, strat Strategy, rec hashSet[string]) (s *reflect.Value, err error) {
	svc, ok := c.services[key]
	if !ok {
		return s, ErrNotRegistered
	}

	if strat == Singleton && svc.value != nil {
		return svc.value, nil
	}

	var val reflect.Value
	if svc.setup != nil {
		v, err := svc.setup(c)
		if err != nil {
			return nil, err
		}
		val = reflect.ValueOf(v)
	} else {
		val = reflect.New(svc.impl)
	}

	elem := val.Elem()

	if rec == nil {
		rec = hashSet[string]{}
	}

	rec.Set(key)

	for i := 0; i < svc.impl.NumField(); i++ {
		field := svc.impl.Field(i)
		stratField := field.Tag.Get("di")
		if stratField == "" {
			continue
		}

		if field.Type.Kind() != reflect.Interface {
			return nil, dependencyError(ErrInvalidDependencyType, field, svc.impl.Name())
		}

		key := getInterfaceKey(field.Type)
		strat, ok := getStrategy(stratField)
		if !ok {
			return nil, dependencyError(ErrInvalidDependencyStrategy, field, svc.impl.Name())
		}

		elemField := elem.Field(i)
		if !elemField.CanSet() {
			return s, dependencyError(ErrDependencyFieldPrivate, field, svc.impl.Name())
		}

		if rec.Has(key) {
			return nil, ErrRecursion
		}

		v, err := c.build(key, strat, rec)
		if err != nil {
			return s, err
		}
		elemField.Set(*v)

		svc.dependencies = append(svc.dependencies, c.services[key])
		if strat == Singleton {
			c.services[key].value = v
		}

		rec.Remove(key)
	}

	if strat == Singleton {
		svc.value = &val
	}

	return &val, nil
}

func MustRegister[TService, TImpl any](c *Container, opt ...RegisterOptions[TImpl]) {
	must(Register[TService, TImpl](c))
}

func MustGet[TService any](c *Container) (s TService) {
	return mustValue(Get[TService](c))
}
