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
		return *new(TInterface), fmt.Errorf("the scope has been closed")
	}
	// resolve
	instance, err := callActivator(ctx, typeOf[TInterface](), options.key)
	if err != nil {
		return *new(TInterface), err
	}
	return instance.(TInterface), nil
}
