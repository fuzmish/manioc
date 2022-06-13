package manioc

import (
	"reflect"
)

// this struct expresses key tuple for registry
type registryKey struct {
	targetType reflect.Type
	key        any
}

// The defaultContext struct implements RegisterContext
type defaultContext struct {
	registry map[registryKey][]activator
	cache    map[any]any
}

// The constructor of defaultContext
func newDefaultContext() *defaultContext {
	return &defaultContext{
		registry: make(map[registryKey][]activator),
		cache:    make(map[any]any),
	}
}

func (ctx *defaultContext) registerActivator(targetType reflect.Type, key any, act activator) error {
	rk := registryKey{targetType: targetType, key: key}
	if _, ok := ctx.registry[rk]; !ok {
		ctx.registry[rk] = make([]activator, 0)
	}
	ctx.registry[rk] = append(ctx.registry[rk], act)
	return nil
}

func (ctx *defaultContext) getActivators(targetType reflect.Type, key any) []activator {
	rk := registryKey{targetType: targetType, key: key}
	if _, ok := ctx.registry[rk]; !ok {
		ctx.registry[rk] = make([]activator, 0)
	}
	return ctx.registry[rk]
}

func (ctx *defaultContext) unregisterActivators(targetType reflect.Type, key any) bool {
	rk := registryKey{targetType: targetType, key: key}
	if ret, ok := ctx.registry[rk]; ok && len(ret) > 0 {
		// clear registration
		ctx.registry[rk] = make([]activator, 0)
		return true
	}
	return false
}

func (ctx *defaultContext) setCache(cacheKey any, value any) {
	ctx.cache[cacheKey] = value
}

func (ctx *defaultContext) getCache(cacheKey any) (any, bool) {
	ret, ok := ctx.cache[cacheKey]
	return ret, ok
}
