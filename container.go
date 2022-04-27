package di

import (
	"sync"
)

type Container interface {
	Put(key string, service Service)
	Get(key string) (Service, bool)
}

type containerImpl struct {
	m sync.Map
}

func NewContainer() Container {
	return &containerImpl{}
}

func (c *containerImpl) Put(key string, service Service) {
	c.m.Store(key, service)
}

func (c *containerImpl) Get(key string) (s Service, ok bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return
	}
	s, ok = v.(Service)
	return
}
