package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSingleton(t *testing.T) {
	c := NewContainer()
	Register[a, struct{}](c)

	s := &singletonService{
		ImplType: reflect.TypeOf(testImpl{}),
		IsBuilt:  false,
	}

	tInstance := s.Build(c)
	assert.NotNil(t, tInstance)
	assert.IsType(t, &testImpl{}, tInstance.Interface())

	impl := tInstance.Interface().(*testImpl)
	assert.NotNil(t, impl.A)

	assert.Equal(t, tInstance, s.Instance)
	assert.True(t, s.IsBuilt)

	tInstance2 := s.Build(c)
	assert.Equal(t, tInstance, tInstance2)
	assert.Same(t, tInstance.Interface(), tInstance2.Interface())
}

func TestBuildTransient(t *testing.T) {
	c := NewContainer()
	Register[a, struct{}](c)

	s := &transientService{
		implType: reflect.TypeOf(testImpl{}),
	}

	tInstance := s.Build(c)
	assert.NotNil(t, tInstance)
	assert.IsType(t, &testImpl{}, tInstance.Interface())

	impl := tInstance.Interface().(*testImpl)
	assert.NotNil(t, impl.A)

	tInstance2 := s.Build(c)
	assert.IsType(t, &testImpl{}, tInstance2.Interface())
	assert.NotEqual(t, tInstance, tInstance2)
	assert.NotSame(t, tInstance.Interface(), tInstance2.Interface())
}
