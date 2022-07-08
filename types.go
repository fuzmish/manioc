package manioc

import (
	"reflect"
)

// CachePolicy is an enumerated type that specifies how the container caches instances.
type CachePolicy int

const (
	// The instance is created only once for the same container, and then it is used repeatedly.
	// This behavior is known as Singleton.
	GlobalCache CachePolicy = iota
	// The instance is created only once for the same scope, and then it is used repeatedly.
	// For different scopes, new instances are created.
	ScopedCache
	// The instance is always newly created and is not cached.
	// This behavior is known as Transient.
	NeverCache
)

// ScopeCacheMode is an enumeration type that configures the behavior of
// the scope with respect to its instance cache.
type ScopeCacheMode int

const (
	// Instance caches are independent across scopes.
	// Even if the parent scope is closed, the child scope will remain open.
	DefaultCacheMode ScopeCacheMode = iota
	// When a child scope is opened, it inherits the instance cache of the parent scope.
	// When the parent scope is closed, the child scopes are also automatically closed.
	InheritCacheMode
	// The parent and child scopes share the instance cache.
	// When the parent scope is closed, the child scopes are also automatically closed.
	SyncCacheMode
)

type registryKey struct {
	serviceType reflect.Type
	serviceKey  any
}

type activator interface {
	activate(ctx resolveContext) (any, error)
}

type resolveContext interface {
	resolve(key registryKey) (any, error)
	setCache(key any, value any, policy CachePolicy)
	getCache(key any, policy CachePolicy) (any, bool)
}

type registerContext interface {
	register(key registryKey, entry activator) error
	isRegistered(key registryKey) bool
	unregister(key registryKey) bool
}

// Scope is an interface that expresses the cache scope of a container.
type Scope interface {
	getResolveContext() resolveContext
	createScope(mode ScopeCacheMode) (Scope, func())
	closeScope()
}

// Container is an interface for storing dependencies.
// It also works as a global scope.
type Container interface {
	Scope
	getRegisterContext() registerContext
}
