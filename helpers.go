package manioc

func NewContainer() Container {
	return newDefaultContainer()
}

func RegisterSingleton[TInterface any, TImplementation any](opts ...RegisterOption) error {
	return Register[TInterface, TImplementation](append(opts, WithCachePolicy(GlobalCache))...)
}

func RegisterScoped[TInterface any, TImplementation any](opts ...RegisterOption) error {
	return Register[TInterface, TImplementation](append(opts, WithCachePolicy(ScopedCache))...)
}

func RegisterTransient[TInterface any, TImplementation any](opts ...RegisterOption) error {
	return Register[TInterface, TImplementation](append(opts, WithCachePolicy(NeverCache))...)
}

func RegisterSingletonConstructor[TInterface any, TConstructor any](ctor TConstructor, opts ...RegisterOption) error {
	return RegisterConstructor[TInterface](ctor, append(opts, WithCachePolicy(GlobalCache))...)
}

func RegisterScopedConstructor[TInterface any, TConstructor any](ctor TConstructor, opts ...RegisterOption) error {
	return RegisterConstructor[TInterface](ctor, append(opts, WithCachePolicy(ScopedCache))...)
}

func RegisterTransientConstructor[TInterface any, TConstructor any](ctor TConstructor, opts ...RegisterOption) error {
	return RegisterConstructor[TInterface](ctor, append(opts, WithCachePolicy(NeverCache))...)
}

func ResolveMany[TInterface any](opts ...ResolveOption) ([]TInterface, error) {
	return Resolve[[]TInterface](opts...)
}
