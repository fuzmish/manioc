package manioc

import (
	"errors"
	"reflect"
)

type defaultContext struct {
	registry    map[registryKey][]activator
	globalCache map[any]any
	scopedCache map[any]any
}

func newDefaultContext() *defaultContext {
	return &defaultContext{
		registry:    make(map[registryKey][]activator),
		globalCache: make(map[any]any),
		scopedCache: make(map[any]any),
	}
}

func (c *defaultContext) register(key registryKey, entry activator) error {
	if _, ok := c.registry[key]; !ok {
		c.registry[key] = make([]activator, 0)
	}
	c.registry[key] = append(c.registry[key], entry)
	return nil
}

func (c *defaultContext) setCache(key any, value any, policy CachePolicy) {
	switch policy {
	case GlobalCache:
		c.globalCache[key] = value
	case ScopedCache:
		c.scopedCache[key] = value
	case NeverCache:
		break
	}
}

func (c *defaultContext) getCache(key any, policy CachePolicy) (any, bool) {
	switch policy {
	case GlobalCache:
		if value, ok := c.globalCache[key]; ok {
			return value, true
		}
	case ScopedCache:
		if value, ok := c.scopedCache[key]; ok {
			return value, true
		}
	case NeverCache:
		break
	}
	return nil, false
}

func (c *defaultContext) resolveAll(key registryKey) (any, error) {
	tkey := registryKey{serviceType: key.serviceType.Elem(), serviceKey: key.serviceKey}
	entries, ok := c.registry[tkey]
	num := len(entries)
	if !ok || num == 0 {
		return nil, errors.New("no registration found")
	}
	// resolve all
	instances := reflect.MakeSlice(key.serviceType, num, num)
	for i, entry := range entries {
		instance, err := entry.activate(c)
		if err != nil {
			return nil, err
		}
		instances.Index(i).Set(reflect.ValueOf(instance))
	}
	return instances.Interface(), nil
}

func (c *defaultContext) resolve(key registryKey) (any, error) {
	// look up entry with key
	entries, ok := c.registry[key]
	if !ok {
		// if service type is []T, look up with T
		if key.serviceType.Kind() == reflect.Slice {
			return c.resolveAll(key)
		}
		return nil, errors.New("no registration found")
	}
	// resolve one
	if len(entries) == 0 {
		return nil, errors.New("no registration found")
	}
	if len(entries) > 1 {
		return nil, errors.New("multiple registration found")
	}
	return entries[0].activate(c)
}

func (c *defaultContext) isRegistered(key registryKey) bool {
	entries, ok := c.registry[key]
	return ok && len(entries) > 0
}

func (c *defaultContext) unregister(key registryKey) bool {
	if ret, ok := c.registry[key]; ok && len(ret) > 0 {
		c.registry[key] = make([]activator, 0)
		return true
	}
	return false
}
