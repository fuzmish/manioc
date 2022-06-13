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

// Activator is a function that creates a service instance on given context.
type activator func(ctx resolveContext) (any, error)

// resolveContext is an interface required for the activator
// to obtain registration information from the container to resolve dependencies.
type resolveContext interface {
	getActivators(targetType reflect.Type, key any) []activator
	setCache(cacheKey any, value any)
	getCache(cacheKey any) (any, bool)
}

// registerContext is an interface required to register dependency information
// and the associated activator into the container.
type registerContext interface {
	resolveContext
	registerActivator(targetType reflect.Type, key any, act activator) error
	unregisterActivators(targetType reflect.Type, key any) bool
}

// Scope is an interface that expresses the cache scope of a container.
type Scope interface {
	getResolveContext() resolveContext
	createScope() (Scope, func())
}

// Container is an interface for storing dependencies.
// It also works as a global scope.
type Container interface {
	Scope
	getRegisterContext() registerContext
}
