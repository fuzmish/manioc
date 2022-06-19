package manioc_container_and_scope_test

import (
	"testing"

	"github.com/fuzmish/manioc"
	"github.com/stretchr/testify/assert"
)

type IMyService interface {
	doSomething()
}

// MyService implements IMyService
type MyService struct {
	Value int
}

func (s *MyService) doSomething() {}

func Test_ContainerAndScope(t *testing.T) {
	assert := assert.New(t)

	// setup container
	ctr := manioc.NewContainer()
	assert.Nil(manioc.Register[IMyService, MyService](manioc.WithContainer(ctr)))

	// open scope
	scope, cleanup := manioc.OpenScope(manioc.WithParentScope(ctr))

	// resolve within the scope
	ret, err := manioc.Resolve[IMyService](manioc.WithScope(scope))
	assert.NotNil(ret)
	assert.Nil(err)

	// close scope
	cleanup()

	// after scope is closed, any resolution within that scope will fail
	_, err = manioc.Resolve[IMyService](manioc.WithScope(scope))
	assert.Error(err)

	// the parent scope (in this case, ctr) is still available
	ret, err = manioc.Resolve[IMyService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)
}

func Test_Scope_WithDefaultCacheMode(t *testing.T) {
	assert := assert.New(t)

	// setup container
	ctr := manioc.NewContainer()

	// register IMyService with ScopedCache policy
	assert.Nil(manioc.RegisterScoped[IMyService, MyService](manioc.WithContainer(ctr)))

	// resolve in parent scope
	ret := manioc.MustResolve[IMyService](manioc.WithScope(ctr))

	// open scope, without inherit/sync instance caches
	scope, cleanup := manioc.OpenScope(
		manioc.WithParentScope(ctr),
		// manioc.WithCacheMode(manioc.DefaultCacheMode),
	)
	defer cleanup()
	retScoped := manioc.MustResolve[IMyService](manioc.WithScope(scope))

	// in DefaultCacheMode, instance caches are independent across scopes
	assert.NotSame(ret, retScoped)
}

func Test_Scope_WithInheritCacheMode(t *testing.T) {
	assert := assert.New(t)

	// setup container
	ctr := manioc.NewContainer()

	// register IMyService with ScopedCache policy
	assert.Nil(manioc.RegisterScoped[IMyService, MyService](manioc.WithContainer(ctr)))

	// resolve in parent scope
	ret := manioc.MustResolve[IMyService](manioc.WithScope(ctr))

	// open scope, inherit instance caches
	scope, cleanup := manioc.OpenScope(
		manioc.WithParentScope(ctr),
		manioc.WithCacheMode(manioc.InheritCacheMode),
	)
	defer cleanup()
	retScoped := manioc.MustResolve[IMyService](manioc.WithScope(scope))

	// since the cache is inherited, the resolution results will match
	assert.Same(ret, retScoped)
}

func Test_Scope_WithSyncCacheMode(t *testing.T) {
	assert := assert.New(t)

	// setup container
	ctr := manioc.NewContainer()

	// register IMyService with ScopedCache policy
	assert.Nil(manioc.RegisterScoped[IMyService, MyService](manioc.WithContainer(ctr)))
	assert.Nil(manioc.RegisterScoped[IMyService, MyService](
		manioc.WithContainer(ctr),
		manioc.WithRegisterKey("another"),
	))

	// resolve in parent scope
	ret := manioc.MustResolve[IMyService](manioc.WithScope(ctr))

	// open scope, sync instance caches
	scope, cleanup := manioc.OpenScope(
		manioc.WithParentScope(ctr),
		manioc.WithCacheMode(manioc.SyncCacheMode),
	)
	defer cleanup()
	retScoped := manioc.MustResolve[IMyService](manioc.WithScope(scope))

	// since the cache is synced, the resolution results will match
	assert.Same(ret, retScoped)

	// resolve in child scope
	ret2Scoped := manioc.MustResolve[IMyService](
		manioc.WithScope(scope),
		manioc.WithResolveKey("another"),
	)

	// since the cache is synced, the resolution results will match
	ret2 := manioc.MustResolve[IMyService](
		manioc.WithScope(ctr),
		manioc.WithResolveKey("another"),
	)
	assert.Same(ret2Scoped, ret2)
}

func Test_NestedScope_WithDefaultCacheMode(t *testing.T) {
	assert := assert.New(t)

	// setup container
	ctr := manioc.NewContainer()
	assert.Nil(manioc.RegisterScoped[IMyService, MyService](manioc.WithContainer(ctr)))

	// open parent scope
	scope, cleanup := manioc.OpenScope(manioc.WithParentScope(ctr))

	// open child scope
	childScope, _ := manioc.OpenScope(
		manioc.WithParentScope(scope),
		// manioc.WithCacheMode(manioc.DefaultCacheMode),
	)

	// close parent
	cleanup()
	_, err := manioc.Resolve[IMyService](manioc.WithScope(scope))
	assert.Error(err)

	// in DefaultCacheMode, the child scope remains open even if the parent scope is closed
	ret, err := manioc.Resolve[IMyService](manioc.WithScope(childScope))
	assert.NotNil(ret)
	assert.Nil(err)
}

func Test_NestedScope_WithInheritCacheMode(t *testing.T) {
	assert := assert.New(t)

	// setup container
	ctr := manioc.NewContainer()
	assert.Nil(manioc.RegisterScoped[IMyService, MyService](manioc.WithContainer(ctr)))

	// open parent scope
	scope, cleanup := manioc.OpenScope(manioc.WithParentScope(ctr))

	// open child scope
	childScope, _ := manioc.OpenScope(
		manioc.WithParentScope(scope),
		manioc.WithCacheMode(manioc.InheritCacheMode),
	)

	// close parent
	cleanup()
	_, err := manioc.Resolve[IMyService](manioc.WithScope(scope))
	assert.Error(err)

	// if you using InheritCacheMode, when the parent scope is closed,
	// then automatically the child scope is also closed.
	_, err = manioc.Resolve[IMyService](manioc.WithScope(childScope))
	assert.Error(err)
}

func Test_NestedScope_WithSyncCacheMode(t *testing.T) {
	assert := assert.New(t)

	// setup container
	ctr := manioc.NewContainer()
	assert.Nil(manioc.RegisterScoped[IMyService, MyService](manioc.WithContainer(ctr)))

	// open parent scope
	scope, cleanup := manioc.OpenScope(manioc.WithParentScope(ctr))

	// open child scope
	childScope, _ := manioc.OpenScope(
		manioc.WithParentScope(scope),
		manioc.WithCacheMode(manioc.SyncCacheMode),
	)

	// close parent
	cleanup()
	_, err := manioc.Resolve[IMyService](manioc.WithScope(scope))
	assert.Error(err)

	// if you using InheritCacheMode, when the parent scope is closed,
	// then automatically the child scope is also closed.
	_, err = manioc.Resolve[IMyService](manioc.WithScope(childScope))
	assert.Error(err)
}
