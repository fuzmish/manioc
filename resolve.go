package manioc

import (
	"errors"
)

func Resolve[TInterface any](opts ...ResolveOption) (TInterface, error) {
	// merge options
	options := &resolveOptions{
		scope: globalContainer,
		key:   nil,
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	// get context
	ctx := options.scope.getResolveContext()
	if ctx == nil {
		return *new(TInterface), errors.New("the scope has been closed")
	}
	// resolve
	instance, err := ctx.resolve(registryKey{
		serviceType: typeof[TInterface](),
		serviceKey:  options.key,
	})
	if err != nil {
		return *new(TInterface), err
	}
	//nolint:forcetypeassert
	return instance.(TInterface), nil
}
