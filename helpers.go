package manioc

func NewContainer() Container {
	return newDefaultContainer()
}

func OpenScope(opts ...OpenScopeOption) (Scope, func()) {
	return globalContainer.OpenScope(opts...)
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

func MustResolve[TInterface any](opts ...ResolveOption) TInterface {
	ret, err := Resolve[TInterface](opts...)
	if err != nil {
		panic(err)
	}
	return ret
}

func ResolveMany[TInterface any](opts ...ResolveOption) ([]TInterface, error) {
	return Resolve[[]TInterface](opts...)
}

func MustResolveMany[TInterface any](opts ...ResolveOption) []TInterface {
	return MustResolve[[]TInterface](opts...)
}

func MustResolveInstance[T any](instance T, opts ...ResolveOption) T {
	ret, err := ResolveInstance(instance, opts...)
	if err != nil {
		panic(err)
	}
	return ret
}

func MustResolveFunction[T any, TFunction any](fun TFunction, opts ...ResolveOption) T {
	ret, err := ResolveFunction[T](fun, opts...)
	if err != nil {
		panic(err)
	}
	return ret
}
