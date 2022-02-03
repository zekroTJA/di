package di

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

//    /--> B --> C --> D
//   /     \--> D
//  A
//   \
//    \--> C --> D

type ServiceA interface {
	DoA()
}

type ServiceB interface {
	DoB()
}

type ServiceC interface {
	DoC()
}

type ServiceD interface {
	DoD()
}

type ImplA struct {
	B ServiceB `di:"singleton"`
	C ServiceC `di:"singleton"`
}

type ImplB struct {
	C ServiceC `di:"singleton"`
	D ServiceD `di:"transistent"`
}

type ImplC struct {
	D ServiceD `di:"transistent"`
}

type ImplD struct {
	some         string
	teardownFunc func() error
}

func (*ImplA) DoA() {}
func (*ImplB) DoB() {}
func (*ImplC) DoC() {}
func (*ImplD) DoD() {}

func (t *ImplD) Teardown() error {
	return t.teardownFunc()
}

func TestGet_General(t *testing.T) {
	c := NewContainer()

	assert.Nil(t, Register[ServiceA, ImplA](c))
	assert.Nil(t, Register[ServiceB, ImplB](c))
	assert.Nil(t, Register[ServiceC, ImplC](c))
	assert.Nil(t, Register[ServiceD](c, RegisterOptions[ImplD]{
		Setup: func(c *Container) (*ImplD, error) {
			return &ImplD{some: "string"}, nil
		},
	}))

	MustGet[ServiceA](c)

	v, err := Get[ServiceD](c)
	assert.Nil(t, err)
	v.DoD()
	assert.Equal(t, "string", v.(*ImplD).some)
	assert.Same(t,
		MustGet[ServiceA](c).(*ImplA).C,
		MustGet[ServiceB](c).(*ImplB).C)
	assert.NotSame(t,
		MustGet[ServiceC](c).(*ImplC).D,
		MustGet[ServiceB](c).(*ImplB).D)
}

func TestTeardown_Func(t *testing.T) {
	c := NewContainer()

	var tdACalled, tdBCalled, tdCCalled uint32
	var cursor uint32

	assert.Nil(t, Register[ServiceA](c, RegisterOptions[ImplA]{
		Teardown: func(c *Container) error {
			tdACalled = atomic.LoadUint32(&cursor)
			atomic.AddUint32(&cursor, 1)
			return nil
		},
	}))
	assert.Nil(t, Register[ServiceB](c, RegisterOptions[ImplB]{
		Teardown: func(c *Container) error {
			tdBCalled = atomic.LoadUint32(&cursor)
			atomic.AddUint32(&cursor, 1)
			return nil
		},
	}))
	assert.Nil(t, Register[ServiceC](c, RegisterOptions[ImplC]{
		Teardown: func(c *Container) error {
			tdCCalled = atomic.LoadUint32(&cursor)
			atomic.AddUint32(&cursor, 1)
			return nil
		},
	}))
	assert.Nil(t, Register[ServiceD](c, RegisterOptions[ImplD]{
		Setup: func(c *Container) (*ImplD, error) {
			return &ImplD{some: "string"}, nil
		},
	}))

	MustGet[ServiceA](c)

	assert.Nil(t, c.Teardown())

	assert.Equal(t, uint32(0), tdCCalled)
	assert.Equal(t, uint32(1), tdBCalled)
	assert.Equal(t, uint32(2), tdACalled)
}

func TestTeardown_Impl(t *testing.T) {
	c := NewContainer()

	assert.Nil(t, Register[ServiceD](c, RegisterOptions[ImplD]{
		Setup: func(c *Container) (*ImplD, error) {
			return &ImplD{some: "string"}, nil
		},
	}))

	var teardownCalled bool

	d := MustGet[ServiceD](c)
	d.(*ImplD).teardownFunc = func() error {
		teardownCalled = true
		return nil
	}

	assert.Nil(t, c.Teardown())

	assert.True(t, teardownCalled)
}
