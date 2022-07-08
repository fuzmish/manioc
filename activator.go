package manioc

import (
	"reflect"
	"unsafe"
)

type implementationActivator struct {
	implementationType reflect.Type
}

func (e *implementationActivator) activate(ctx resolveContext) (any, error) {
	var instance any
	if e.implementationType.Kind() == reflect.Pointer {
		instance = reflect.New(e.implementationType.Elem()).Interface()
	} else {
		instance = reflect.New(e.implementationType).Elem().Interface()
	}
	return instance, nil
}

type instanceActivator struct {
	instance any
}

func (e *instanceActivator) activate(ctx resolveContext) (any, error) {
	return e.instance, nil
}

type constructorActivator struct {
	constructor any
}

func (e *constructorActivator) activate(ctx resolveContext) (any, error) {
	tFn := reflect.TypeOf(e.constructor)
	vFn := reflect.ValueOf(e.constructor)
	numArgs := tFn.NumIn()
	tFnArgs := make([]reflect.Type, numArgs)
	for i := 0; i < numArgs; i++ {
		tFnArgs[i] = tFn.In(i)
	}
	// constructor injection
	args := make([]reflect.Value, numArgs)
	for idx := 0; idx < numArgs; idx++ {
		instance, err := ctx.resolve(registryKey{
			serviceType: tFnArgs[idx],
			serviceKey:  nil, /* no key is available for constructor injection */
		})
		if err != nil {
			return nil, err
		}
		args[idx] = reflect.ValueOf(instance)
	}
	ret := vFn.Call(args)
	// check error value
	if len(ret) == 2 && ret[1].IsValid() && !ret[1].IsNil() {
		//nolint:forcetypeassert
		err := ret[1].Interface().(error)
		if err != nil {
			return nil, err
		}
	}
	instance := ret[0].Interface()
	return instance, nil
}

type fieldInjectionActivator struct {
	baseActivator activator
}

func (e *fieldInjectionActivator) activate(ctx resolveContext) (any, error) {
	instance, err := e.baseActivator.activate(ctx)
	if err != nil {
		return nil, err
	}
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	// skip if instance is not a struct
	if val.Kind() != reflect.Struct {
		return instance, nil
	}
	// field injection
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := field.Type()
		info, err := parseTag(t.Field(i).Tag)
		if err != nil {
			return nil, err
		}
		if !info.inject {
			continue
		}
		if !field.CanSet() {
			// accessing unexported fields
			// cf. https://stackoverflow.com/a/43918797
			field = reflect.NewAt(fieldType, unsafe.Pointer(field.UnsafeAddr())).Elem()
		}
		instance, err := ctx.resolve(registryKey{
			serviceType: fieldType,
			serviceKey:  info.key,
		})
		if err != nil {
			return nil, err
		}
		field.Set(reflect.ValueOf(instance))
	}
	return instance, nil
}

type cacheActivator struct {
	baseActivator activator
	policy        CachePolicy
}

func (e *cacheActivator) activate(ctx resolveContext) (any, error) {
	// check cache
	if instance, ok := ctx.getCache(e, e.policy); ok {
		return instance, nil
	}
	// activate new instance
	instance, err := e.baseActivator.activate(ctx)
	if err != nil {
		return nil, err
	}
	// store
	ctx.setCache(e, instance, e.policy)
	return instance, nil
}
