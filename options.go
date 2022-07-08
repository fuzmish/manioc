package manioc

import (
	"fmt"
)

// options for Register
type registerOptions struct {
	container Container
	key       any
	policy    CachePolicy
}

type RegisterOption interface {
	apply(*registerOptions)
}

// WithContainer

type withContainer struct{ container Container }

func (opt *withContainer) apply(options *registerOptions) {
	options.container = opt.container
}

func WithContainer(container Container) RegisterOption {
	return &withContainer{container: container}
}

// WithRegisterKey

type withRegisterKey struct{ key any }

func (opt *withRegisterKey) apply(options *registerOptions) {
	options.key = opt.key
}

func WithRegisterKey(key any) RegisterOption {
	return &withRegisterKey{key: key}
}

// WithLifeTime

type withCachePolicy struct{ policy CachePolicy }

func (opt *withCachePolicy) apply(options *registerOptions) {
	options.policy = opt.policy
}

func WithCachePolicy(policy CachePolicy) RegisterOption {
	switch policy {
	case GlobalCache, ScopedCache, NeverCache:
		return &withCachePolicy{policy: policy}
	default:
		panic(fmt.Errorf("invalid CachePolicy value: `%v`", policy))
	}
}

//
// options for Resolve
//

type resolveOptions struct {
	scope Scope
	key   any
}

type ResolveOption interface {
	apply(*resolveOptions)
}

// WithScope

type withScope struct{ scope Scope }

func (opt *withScope) apply(options *resolveOptions) {
	options.scope = opt.scope
}

func WithScope(scope Scope) ResolveOption {
	return &withScope{scope: scope}
}

// WithResolveKey

type withResolveKey struct{ key any }

func (opt *withResolveKey) apply(options *resolveOptions) {
	options.key = opt.key
}

func WithResolveKey(key any) ResolveOption {
	return &withResolveKey{key: key}
}

//
// options for OpenScope
//

type openScopeOptions struct {
	parent    Scope
	cacheMode ScopeCacheMode
}

type OpenScopeOption interface {
	apply(*openScopeOptions)
}

// WithParentScope

type withParentScope struct{ parent Scope }

func (opt *withParentScope) apply(options *openScopeOptions) {
	options.parent = opt.parent
}

func WithParentScope(parent Scope) OpenScopeOption {
	return &withParentScope{parent: parent}
}

// WithCacheMode

type withCacheMode struct{ cacheMode ScopeCacheMode }

func (opt *withCacheMode) apply(options *openScopeOptions) {
	options.cacheMode = opt.cacheMode
}

func WithCacheMode(cacheMode ScopeCacheMode) OpenScopeOption {
	return &withCacheMode{cacheMode: cacheMode}
}
