package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type a interface{}

type testImpl struct {
	A a
}

func TestRegister(t *testing.T) {
	c := NewContainer()

	// Infers testImpl as type instead
	// of an interface.
	err := Register(c, testImpl{})
	assert.ErrorIs(t, err, ErrNoInterface)

	type myInterface interface{}
	const key = "github.com/zekrotja/di.myInterface"

	err = Register[myInterface](c, testImpl{})
	assert.Nil(t, err)

	v, ok := c.(*containerImpl).m.Load(key)
	assert.True(t, ok, "no service has been registered")

	svc := v.(*Service)
	assert.Equal(t, svc.ImplType, reflect.TypeOf(testImpl{}))
	assert.False(t, svc.IsBuilt)
	assert.Equal(t, svc.Instance, reflect.Value{})
}

func TestGet(t *testing.T) {
	c := NewContainer()

	type myInterface interface{}
	type myOtherInterface interface{}

	impl := testImpl{}
	Register[myInterface](c, impl)

	s, err := Get[myInterface](c)
	assert.Nil(t, err)
	assert.IsType(t, s, &impl)

	// Ensure that the retrieved value is
	// exactly the same instance on an
	// repeated retrieve.
	s2, err := Get[myInterface](c)
	assert.Nil(t, err)
	assert.Same(t, s2, s)

	_, err = Get[myOtherInterface](c)
	assert.ErrorIs(t, err, ErrNotRegistered)

	_, err = Get[struct{}](c)
	assert.ErrorIs(t, err, ErrNoInterface)
}
