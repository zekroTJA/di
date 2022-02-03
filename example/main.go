package main

import (
	"fmt"

	"github.com/zekrotja/di"
)

type A interface {
	Some()
}

type B interface {
	More()
}

type aImpl struct{}

func (aImpl) Some() {
	fmt.Println("some")
}

type bImpl struct {
	A A
	B int
}

func (b bImpl) More() {
	b.A.Some()
}

func main() {
	c := di.NewContainer()
	fmt.Println(di.Register[A](c, aImpl{}))
	fmt.Println(di.Register[B](c, bImpl{}))
	b, _ := di.Get[B](c)
	b.More()
}
