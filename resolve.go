package manioc

import (
	"fmt"
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
		//nolint:gocritic
		return *new(TInterface), fmt.Errorf("the scope has been closed")
	}
	// resolve
	instance, err := callActivator(ctx, typeof[TInterface](), options.key)
	if err != nil {
		//nolint:gocritic
		return *new(TInterface), err
	}
	//nolint:forcetypeassert
	return instance.(TInterface), nil
}
