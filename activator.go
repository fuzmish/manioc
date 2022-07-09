package manioc

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type implementationActivator struct {
	implementationType reflect.Type
}

func newImplementationActivator[TInterface any, TImplementation any]() activator {
	// check type parameter
	tIface := typeof[TInterface]()
	tImpl := typeof[TImplementation]()
	if tIface.Kind() == reflect.Interface && tImpl.Kind() != reflect.Pointer {
		tImpl = reflect.PointerTo(tImpl)
	}
	tImplElm := tImpl
	if tImpl.Kind() == reflect.Pointer {
		tImplElm = tImpl.Elem()
	}
	if tImplElm.Kind() == reflect.Interface {
		panic(fmt.Errorf(
			"TImplementation=`%s` should not be an interface type",
			nameof[TImplementation](),
		))
	}
	if !tImpl.AssignableTo(tIface) {
		panic(fmt.Errorf(
			"TImplementation=`%s` should be assignable to TInterface=`%s`",
			nameof[TImplementation](),
			nameof[TInterface](),
		))
	}
	return &implementationActivator{implementationType: tImpl}
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

func newInstanceActivator(instance any) (activator, error) {
	if !reflect.ValueOf(instance).IsValid() {
		return nil, errors.New("instance is invalid")
	}
	//nolint:exhaustive
	switch reflect.TypeOf(instance).Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Pointer,
		reflect.UnsafePointer,
		reflect.Interface,
		reflect.Slice:
		if reflect.ValueOf(instance).IsNil() {
			return nil, errors.New("instance is nil")
		}
	}
	return &instanceActivator{instance: instance}, nil
}

func (e *instanceActivator) activate(ctx resolveContext) (any, error) {
	return e.instance, nil
}

type constructorActivator struct {
	constructor any
}

func newConstructorActivator[T any, TConstructor any](ctor TConstructor) (activator, error) {
	// check type parameters
	tRet := typeof[T]()
	tCtor := typeof[TConstructor]()
	if tCtor.Kind() != reflect.Func {
		panic(errors.New("the type of TConstructor should be a function"))
	}
	switch tCtor.NumOut() {
	case 1:
		// out[0] should be assignable to T
		if !tCtor.Out(0).AssignableTo(tRet) {
			panic(fmt.Errorf(
				"the return value TConstructor=`%s` should be assignable to T=`%s`",
				nameof[TConstructor](),
				nameof[T](),
			))
		}
	case 2:
		// out[0] should be assignable to T
		if !tCtor.Out(0).AssignableTo(tRet) {
			panic(fmt.Errorf(
				"the first return value TConstructor=`%s` should be assignable to T=`%s`",
				nameof[TConstructor](),
				nameof[T](),
			))
		}
		// out[1] should be an error
		if tCtor.Out(1) != typeof[error]() {
			panic(fmt.Errorf(
				"the second return value of TConstructor=`%s` should be error",
				nameof[TConstructor](),
			))
		}
	default:
		panic(errors.New("unexpected number of return values"))
	}
	// check ctor value
	if !reflect.ValueOf(ctor).IsValid() || reflect.ValueOf(ctor).IsNil() {
		return nil, errors.New("ctor is invalid or nil")
	}
	return &constructorActivator{constructor: ctor}, nil
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
