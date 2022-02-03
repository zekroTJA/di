package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	c := NewContainer()
	Register[a, struct{}](c)

	s := &Service{
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
}
