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

func ensureInterface[T any]() {
	if typeof[T]().Kind() != reflect.Interface {
		panic(fmt.Errorf("EnsureInterface: %s is not an interface", nameof[T]()))
	}
}

func ensureImplements[TInterface any, TImplementation any]() {
	if typeof[TInterface]().Kind() != reflect.Interface {
		panic(fmt.Errorf("EnsureImplements: %s is not an interface", nameof[TInterface]()))
	}
	if typeof[TImplementation]().Kind() == reflect.Interface {
		panic(fmt.Errorf("EnsureImplements: an interface type %s is not allowed for TImplementation",
			nameof[TImplementation]()))
	}
	if !typeof[*TImplementation]().Implements(typeof[TInterface]()) {
		panic(fmt.Errorf("EnsureImplements: %s does not implement %s",
			nameof[TImplementation](),
			nameof[TInterface]()))
	}
}

func ensureFunctionReturnType[TFunc any, TReturn any]() {
	funcType := typeof[TFunc]()
	if funcType.Kind() != reflect.Func {
		panic(fmt.Errorf("EnsureFunctionReturnType: TFunc %v is not a function type",
			nameof[TFunc]()))
	}
	if funcType.NumOut() != 1 {
		panic(fmt.Errorf("EnsureFunctionReturnType: At this moment the function returns multi-value is not supported"))
	}
	if funcType.IsVariadic() {
		panic(fmt.Errorf("EnsureFunctionReturnType: At this moment variadic arg is not supported"))
	}
	out := funcType.Out(0)
	ret := typeof[TReturn]()
	if !out.AssignableTo(ret) {
		panic(fmt.Errorf("EnsureFunctionReturnType: The return type %v of TFunc %v is not assignable to %v",
			out.Name(),
			nameof[TFunc](),
			nameof[TReturn]()))
	}
}
