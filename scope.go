package manioc

type defaultScope struct {
	context     *defaultContext
	childScopes []Scope
}

func (c *defaultScope) getResolveContext() resolveContext {
	if c.context == nil {
		return nil
	}
	return c.context
}

func (c *defaultScope) createScope(mode ScopeCacheMode) (Scope, func()) {
	ret := &defaultScope{
		context: &defaultContext{
			registry:    c.context.registry,
			globalCache: c.context.globalCache,
			scopedCache: make(map[any]any),
		},
		childScopes: make([]Scope, 0),
	}
	if mode == InheritCacheMode {
		// inherit parent cache
		for k, v := range c.context.scopedCache {
			ret.context.scopedCache[k] = v
		}
		// register child scope into parent
		c.childScopes = append(c.childScopes, ret)
	} else if mode == SyncCacheMode {
		// syncrhonize cache
		ret.context.scopedCache = c.context.scopedCache
		// register child scope into parent
		c.childScopes = append(c.childScopes, ret)
	}
	cleanup := func() {
		// after this function is called, this scope is no longer available.
		ret.closeScope()
	}
	return ret, cleanup
}

func (c *defaultScope) closeScope() {
	if c.childScopes != nil {
		for _, scope := range c.childScopes {
			scope.closeScope()
		}
		c.childScopes = nil
		c.context = nil
	}
}

// Create new child scope.
// The second return value is a cleanup function, which you can call it
// to explicitly close the scope. After this function is called,
// the resolution request for the corresponding scope will not work.
func OpenScope(opts ...OpenScopeOption) (Scope, func()) {
	// merge options
	options := &openScopeOptions{
		parent:    globalContainer,
		cacheMode: DefaultCacheMode,
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	// create scoped container
	return options.parent.createScope(options.cacheMode)
}
