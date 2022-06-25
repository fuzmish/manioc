package manioc

import (
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
	ret := ctx.getActivators(typeof[TInterface](), options.key)
	return len(ret) > 0
}

func registerActivator[TInterface any](act activator, opts ...RegisterOption) error {
	// parse options
	options := mergeRegisterOptions(opts)
	// install instance cache activator
	act = createCachedActivator(act, options.policy)
	// register
	ctx := options.container.getRegisterContext()
	//nolint:wrapcheck
	return ctx.registerActivator(typeof[TInterface](), options.key, act)
}

func RegisterConstructor[TInterface any, TConstructor any](ctor TConstructor, opts ...RegisterOption) error {
	// create activator with constructor injection
	activator := createConstructorInjectionActivator[TInterface](ctor)
	// install field injection activator
	activator = createFieldInjectionActivator(activator)
	// register
	return registerActivator[TInterface](activator, opts...)
}

func RegisterInstance[TInterface any](instance TInterface, opts ...RegisterOption) error {
	// verify instance
	if !reflect.ValueOf(instance).IsValid() {
		return fmt.Errorf("cannot register nil as an instance")
	}
	// create singleton instance activator
	activator := createSingletonInstanceActivator(instance)
	// install field injection activator
	activator = createFieldInjectionActivator(activator)
	// overwrite cache policy option with GlobalCache
	opts = append(opts, WithCachePolicy(GlobalCache))
	// register
	return registerActivator[TInterface](activator, opts...)
}

func Register[TInterface any, TImplementation any](opts ...RegisterOption) error {
	// create activator with default constructor
	activator := createImplementationActivator[TInterface, TImplementation]()
	// install field injection activator
	activator = createFieldInjectionActivator(activator)
	// register
	return registerActivator[TInterface](activator, opts...)
}

func Unregister[TInterface any](opts ...RegisterOption) bool {
	options := mergeRegisterOptions(opts)
	ctx := options.container.getRegisterContext()
	return ctx.unregisterActivators(typeof[TInterface](), options.key)
}
