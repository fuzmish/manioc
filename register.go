package manioc

import (
	"errors"
	"fmt"
	"reflect"
)

func mergeRegisterOptions(opts []RegisterOption) *registerOptions {
	options := &registerOptions{
		container: globalContainer,
		key:       nil,
		policy:    NeverCache,
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	return options
}

func IsRegistered[TInterface any](opts ...RegisterOption) bool {
	options := mergeRegisterOptions(opts)
	ctx := options.container.getRegisterContext()
	key := registryKey{serviceType: typeof[TInterface](), serviceKey: options.key}
	return ctx.isRegistered(key)
}

func register(serviceType reflect.Type, activator activator, opts ...RegisterOption) error {
	// parse option
	options := mergeRegisterOptions(opts)
	// get context
	ctx := options.container.getRegisterContext()
	// install field injection activator
	activator = &fieldInjectionActivator{baseActivator: activator}
	// install cache activator
	activator = &cacheActivator{baseActivator: activator, policy: options.policy}
	// register
	key := registryKey{serviceType: serviceType, serviceKey: options.key}
	return ctx.register(key, activator)
}

func RegisterConstructor[TInterface any, TConstructor any](ctor TConstructor, opts ...RegisterOption) error {
	// check type parameters
	tIface := typeof[TInterface]()
	tCtor := typeof[TConstructor]()
	if tCtor.Kind() != reflect.Func {
		panic(errors.New("the type of TConstructor should be a function"))
	}
	switch tCtor.NumOut() {
	case 1:
		// out[0] should be assignable to TInterface
		if !tCtor.Out(0).AssignableTo(tIface) {
			panic(fmt.Errorf(
				"the return value TConstructor=`%s` should be assignable to TInterface=`%s`",
				nameof[TConstructor](),
				nameof[TInterface](),
			))
		}
	case 2:
		// out[0] should be assignable to TInterface
		if !tCtor.Out(0).AssignableTo(tIface) {
			panic(fmt.Errorf(
				"the first return value TConstructor=`%s` should be assignable to TInterface=`%s`",
				nameof[TConstructor](),
				nameof[TInterface](),
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
		return errors.New("the value of ctor is invalid")
	}
	// register
	return register(
		tIface,
		&constructorActivator{constructor: ctor},
		opts...,
	)
}

func RegisterInstance[TInterface any](instance TInterface, opts ...RegisterOption) error {
	// check instance value
	if !reflect.ValueOf(instance).IsValid() {
		return errors.New("instance is nil")
	}
	return register(
		typeof[TInterface](),
		&instanceActivator{instance: instance},
		append(opts, WithCachePolicy(GlobalCache))...,
	)
}

func Register[TInterface any, TImplementation any](opts ...RegisterOption) error {
	tIface := typeof[TInterface]()
	tImpl := typeof[TImplementation]()
	// check type parameter
	if tIface.Kind() == reflect.Interface && tImpl.Kind() != reflect.Pointer {
		tImpl = reflect.PointerTo(tImpl)
	}
	if !tImpl.AssignableTo(tIface) {
		panic(fmt.Errorf(
			"TImplementation=`%s` should be assignable to TInterface=`%s`",
			nameof[TImplementation](),
			nameof[TInterface](),
		))
	}
	return register(
		tIface,
		&implementationActivator{implementationType: tImpl},
		opts...,
	)
}

func Unregister[TInterface any](opts ...RegisterOption) bool {
	options := mergeRegisterOptions(opts)
	ctx := options.container.getRegisterContext()
	key := registryKey{serviceType: typeof[TInterface](), serviceKey: options.key}
	return ctx.unregister(key)
}
