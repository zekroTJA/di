package di

import (
	"fmt"
	"reflect"
)

const (
	ErrSvcNoInterface            = Error("service must be an interface type")
	ErrImplNoStruct              = Error("implementation must be a struct type")
	ErrImplDoesNotImplementSvc   = Error("the registered implementation does not implement the given service")
	ErrInvalidDependencyStrategy = Error("the specified strategy is invalid")
	ErrInvalidDependencyType     = Error("a depdendency must be of type interface")
	ErrNotRegistered             = Error("service has not been registered")
	ErrRecursion                 = Error("recursion in dependency graph")
	ErrStoredTypeCanNotBeCasted  = Error("the stored type can not be casted to the service interface (this should never happen)")
	ErrDependencyFieldPrivate    = Error("the dependency field is private and can thus not be set by the DI system")
)

type Error string

func (t Error) Error() string {
	return string(t)
}

type InnerError struct {
	Inner error
}

func (t InnerError) Error() string {
	return t.Inner.Error()
}

func (t InnerError) Unwrap() error {
	return t.Inner
}

// --- Specific Error Types ---

type DependencyError struct {
	InnerError
	FieldName  string
	StructName string
}

func dependencyError(inner error, field reflect.StructField, structName string) error {
	return DependencyError{
		InnerError: InnerError{
			Inner: inner,
		},
		FieldName:  field.Name,
		StructName: structName,
	}
}

func (t DependencyError) Error() string {
	return fmt.Sprintf("dependency error %s.%s: %s", t.StructName, t.FieldName, t.Inner)
}
