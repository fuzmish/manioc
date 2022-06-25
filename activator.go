package manioc

import (
	"fmt"
	"reflect"
	"unsafe"
)

func callActivator(ctx resolveContext, targetType reflect.Type, key any) (any, error) {
	// if targetType is slice, use all activators
	if targetType.Kind() == reflect.Slice {
		activators := ctx.getActivators(targetType.Elem(), key)
		num := len(activators)
		if num == 0 {
			return nil, fmt.Errorf("no activator found")
		}
		slice := reflect.MakeSlice(targetType, num, num)
		for i, act := range activators {
			instance, err := act(ctx)
			if err != nil {
				return nil, err
			}
			slice.Index(i).Set(reflect.ValueOf(instance))
		}
		return slice.Interface(), nil
	}
	// otherwise, use the first activator
	activators := ctx.getActivators(targetType, key)
	if len(activators) == 0 {
		return nil, fmt.Errorf("no activator found")
	}
	if len(activators) > 1 {
		return nil, fmt.Errorf("there are multiple activators and we cannot choose which one to use")
	}
	// instantiate
	activator := activators[0]
	instance, err := activator(ctx)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func injectToFields(ctx resolveContext, instance any) error {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("field injection is not allowed for non-struct value")
	}
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := field.Type()
		info, err := parseTag(t.Field(i).Tag)
		if err != nil {
			return err
		}
		if !info.inject {
			continue
		}
		if !field.CanSet() {
			// accessing unexported fields
			// cf. https://stackoverflow.com/a/43918797
			field = reflect.NewAt(fieldType, unsafe.Pointer(field.UnsafeAddr())).Elem()
		}
		instance, err := callActivator(ctx, fieldType, info.key)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(instance))
	}
	return nil
}

func createImplementationActivator[TInterface any, TImplementation any]() activator {
	// check type parameters;
	// It would be great if this check were performed statically,
	// but this is not possible in current Golang.
	// Therefore we will instead perform the inspection here at runtime.
	ensureImplements[TInterface, TImplementation]()
	// create activator
	return func(ctx resolveContext) (any, error) { return new(TImplementation), nil }
}

func createConstructorInjectionActivator[TInterface any, TConstructor any](ctor TConstructor) activator {
	// check type parameters;
	// It would be great if this check were performed statically,
	// but this is not possible in current Golang.
	// Therefore we will instead perform the inspection here at runtime.
	ensureFunctionReturnType[TConstructor, TInterface]()
	// cache constructor reflect info
	tFn := typeof[TConstructor]()
	vFn := reflect.ValueOf(ctor)
	numArgs := tFn.NumIn()
	tFnArgs := make([]reflect.Type, numArgs)
	for i := 0; i < numArgs; i++ {
		tFnArgs[i] = tFn.In(i)
	}
	// activator with constructor injection
	return func(ctx resolveContext) (any, error) {
		// constructor injection
		args := make([]reflect.Value, numArgs)
		for i := 0; i < numArgs; i++ {
			instance, err := callActivator(ctx, tFnArgs[i], nil /* no key is available for constructor injection */)
			if err != nil {
				return nil, err
			}
			args[i] = reflect.ValueOf(instance)
		}
		ret := vFn.Call(args)
		instance := ret[0].Interface()
		return instance, nil
	}
}

func createSingletonInstanceActivator(instance any) activator {
	// create identity activator
	return func(ctx resolveContext) (any, error) {
		return instance, nil
	}
}

func createFieldInjectionActivator(baseActivator activator) activator {
	return func(ctx resolveContext) (any, error) {
		instance, err := baseActivator(ctx)
		if err != nil {
			return nil, err
		}
		if err := injectToFields(ctx, instance); err != nil {
			return nil, err
		}
		return instance, nil
	}
}

func createCachedActivator(baseActivator activator, policy CachePolicy) activator {
	if policy == NeverCache {
		// for NeverCache policy, instance cache wrapper is not needed
		return baseActivator
	}
	if policy == ScopedCache || policy == GlobalCache {
		isGlobal := policy == GlobalCache
		var act activator
		act = func(ctx resolveContext) (any, error) {
			if _, ok := ctx.getCache(&act, isGlobal); !ok {
				ret, err := baseActivator(ctx)
				if err != nil {
					return nil, err
				}
				ctx.setCache(&act, ret, isGlobal)
				return ret, nil
			}
			ret, _ := ctx.getCache(&act, isGlobal)
			return ret, nil
		}
		return act
	}
	panic(fmt.Errorf("invalid CachePolicy value: %v", policy))
}
