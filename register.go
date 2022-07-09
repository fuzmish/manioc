package manioc

import (
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

func IsRegistered[T any](opts ...RegisterOption) bool {
	options := mergeRegisterOptions(opts)
	ctx := options.container.getRegisterContext()
	key := registryKey{serviceType: typeof[T](), serviceKey: options.key}
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

func RegisterConstructor[T any, TConstructor any](ctor TConstructor, opts ...RegisterOption) error {
	activator, err := newConstructorActivator[T](ctor)
	if err != nil {
		return err
	}
	return register(typeof[T](), activator, opts...)
}

func RegisterInstance[T any](instance T, opts ...RegisterOption) error {
	activator, err := newInstanceActivator(instance)
	if err != nil {
		return err
	}
	// override cache policy
	opts = append(opts, WithCachePolicy(GlobalCache))
	return register(typeof[T](), activator, opts...)
}

func Register[TInterface any, TImplementation any](opts ...RegisterOption) error {
	activator := newImplementationActivator[TInterface, TImplementation]()
	return register(typeof[TInterface](), activator, opts...)
}

func Unregister[T any](opts ...RegisterOption) bool {
	options := mergeRegisterOptions(opts)
	ctx := options.container.getRegisterContext()
	key := registryKey{serviceType: typeof[T](), serviceKey: options.key}
	return ctx.unregister(key)
}
