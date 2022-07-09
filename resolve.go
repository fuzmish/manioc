package manioc

import (
	"errors"
)

func mergeResolveOptions(opts []ResolveOption) *resolveOptions {
	options := &resolveOptions{
		scope: globalContainer,
		key:   nil,
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	return options
}

func Resolve[T any](opts ...ResolveOption) (T, error) {
	// parse option
	options := mergeResolveOptions(opts)
	// get context
	ctx := options.scope.getResolveContext()
	if ctx == nil {
		return *new(T), errors.New("the scope has been closed")
	}
	// resolve
	instance, err := ctx.resolve(registryKey{
		serviceType: typeof[T](),
		serviceKey:  options.key,
	})
	if err != nil {
		return *new(T), err
	}
	//nolint:forcetypeassert
	return instance.(T), nil
}

func directResolve(activator activator, opts ...ResolveOption) (any, error) {
	// parse option
	options := mergeResolveOptions(opts)
	// get context
	ctx := options.scope.getResolveContext()
	// install field injection activator
	activator = &fieldInjectionActivator{baseActivator: activator}
	// resolve
	ret, err := activator.activate(ctx)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func ResolveInstance[T any](instance T, opts ...ResolveOption) (T, error) {
	// using instance activator
	activator, err := newInstanceActivator(instance)
	if err != nil {
		return *new(T), err
	}
	ret, err := directResolve(activator, opts...)
	if err != nil {
		return *new(T), err
	}
	//nolint:forcetypeassert
	return ret.(T), nil
}

func ResolveFunction[T any, TFunction any](fun TFunction, opts ...ResolveOption) (T, error) {
	// using constructor activator
	activator, err := newConstructorActivator[T](fun)
	if err != nil {
		return *new(T), err
	}
	ret, err := directResolve(activator, opts...)
	if err != nil {
		return *new(T), err
	}
	//nolint:forcetypeassert
	return ret.(T), nil
}
