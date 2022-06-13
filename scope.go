package manioc

type defaultScope struct {
	context *defaultContext
}

func (c *defaultScope) getResolveContext() resolveContext {
	if c.context == nil {
		return nil
	}
	return c.context
}

func (c *defaultScope) createScope() (Scope, func()) {
	ret := &defaultScope{
		context: &defaultContext{
			registry: c.context.registry,
			cache:    make(map[any]any),
		},
	}
	cleanup := func() {
		// after this function is called, this scope is no longer available.
		ret.context = nil
	}
	return ret, cleanup
}

// Create new child scope.
// The created scope inherits dependency registration,
// but does not inherit the instance cache.
// The second return value is a cleanup function, which you can call it
// to explicitly close the scope. After this function is called,
// the resolution request for the corresponding scope will not work.
func OpenScope(opts ...OpenScopeOption) (Scope, func()) {
	// merge options
	options := &openScopeOptions{
		parent: globalContainer,
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	// create scoped container
	return options.parent.createScope()
}
