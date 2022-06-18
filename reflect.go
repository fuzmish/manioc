package manioc

import (
	"fmt"
	"reflect"
)

func typeof[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func nameof[T any]() string {
	return typeof[T]().Name()
}

func ensureInterface[T any]() error {
	if typeof[T]().Kind() == reflect.Interface {
		return nil
	}
	return fmt.Errorf("EnsureInterface: %s is not an interface", nameof[T]())
}

func ensureImplements[TInterface any, TImplementation any]() error {
	if err := ensureInterface[TInterface](); err != nil {
		return err
	}
	if err := ensureInterface[TImplementation](); err == nil {
		return fmt.Errorf("EnsureImplements: an interface type %s is not allowed for TImplementation",
			nameof[TImplementation]())
	}
	if typeof[*TImplementation]().Implements(typeof[TInterface]()) {
		return nil
	}
	return fmt.Errorf("EnsureImplements: %s does not implement %s",
		nameof[TImplementation](),
		nameof[TInterface]())
}

func ensureFunctionReturnType[TFunc any, TReturn any]() error {
	funcType := typeof[TFunc]()
	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("EnsureFunctionReturnType: TFunc %v is not a function type",
			nameof[TFunc]())
	}
	if funcType.NumOut() != 1 {
		return fmt.Errorf("EnsureFunctionReturnType: At this moment the function returns multi-value is not supported")
	}
	if funcType.IsVariadic() {
		return fmt.Errorf("EnsureFunctionReturnType: At this moment variadic arg is not supported")
	}
	out := funcType.Out(0)
	ret := typeof[TReturn]()
	if !out.AssignableTo(ret) {
		return fmt.Errorf("EnsureFunctionReturnType: The return type %v of TFunc %v is not assignable to %v",
			out.Name(),
			nameof[TFunc](),
			nameof[TReturn]())
	}
	return nil
}
