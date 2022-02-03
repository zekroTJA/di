package di

import (
	"fmt"
	"reflect"
)

func getInterfaceType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func getInterfaceKey(t reflect.Type) string {
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}
