package manioc

import (
	"fmt"
	"reflect"
)

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func nameOf[T any]() string {
	return typeOf[T]().Name()
}

func ensureInterface[T any]() error {
	if typeOf[T]().Kind() == reflect.Interface {
		return nil
	}
	return fmt.Errorf("EnsureInterface: %s is not an interface", nameOf[T]())
}

func ensureImplements[TInterface any, TImplementation any]() error {
	if err := ensureInterface[TInterface](); err != nil {
		return err
	}
	if err := ensureInterface[TImplementation](); err == nil {
		return fmt.Errorf("EnsureImplements: an interface type %s is not allowed for TImplementation",
			nameOf[TImplementation]())
	}
	if typeOf[*TImplementation]().Implements(typeOf[TInterface]()) {
		return nil
	}
	return fmt.Errorf("EnsureImplements: %s does not implement %s",
		nameOf[TImplementation](),
		nameOf[TInterface]())
}

func ensureFunctionReturnType[TFunc any, TReturn any]() error {
	fn := typeOf[TFunc]()
	if fn.Kind() != reflect.Func {
		return fmt.Errorf("EnsureFunctionReturnType: TFunc %v is not a function type",
			nameOf[TFunc]())
	}
	if fn.NumOut() != 1 {
		return fmt.Errorf("EnsureFunctionReturnType: At this moment the function returns multi-value is not supported")
	}
	if fn.IsVariadic() {
		return fmt.Errorf("EnsureFunctionReturnType: At this moment variadic arg is not supported")
	}
	out := fn.Out(0)
	ret := typeOf[TReturn]()
	if !out.AssignableTo(ret) {
		return fmt.Errorf("EnsureFunctionReturnType: The return type %v of TFunc %v is not assignable to %v",
			out.Name(),
			nameOf[TFunc](),
			nameOf[TReturn]())
	}
	return nil
}
