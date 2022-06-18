package manioc

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
	options := mergeRegisterOptions(opts)
	ctx := options.container.getRegisterContext()
	//nolint:wrapcheck
	return ctx.registerActivator(typeof[TInterface](), options.key, createCachedActivator(act, options.policy))
}

func RegisterConstructor[TInterface any, TConstructor any](ctor TConstructor, opts ...RegisterOption) error {
	return registerActivator[TInterface](createConstructorActivator[TInterface](ctor), opts...)
}

func RegisterInstance[TInterface any](instance TInterface, opts ...RegisterOption) error {
	return registerActivator[TInterface](createSingletonInstanceActivator(instance),
		append(opts, WithCachePolicy(GlobalCache))...)
}

func Register[TInterface any, TImplementation any](opts ...RegisterOption) error {
	return registerActivator[TInterface](createImplementationActivator[TInterface, TImplementation](), opts...)
}

func Unregister[TInterface any](opts ...RegisterOption) bool {
	options := mergeRegisterOptions(opts)
	ctx := options.container.getRegisterContext()
	return ctx.unregisterActivators(typeof[TInterface](), options.key)
}
