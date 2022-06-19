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

func Test_NestedScope(t *testing.T) {
	assert := assert.New(t)

	// setup container
	ctr := manioc.NewContainer()

	// register IMyService with ScopedCache policy
	assert.Nil(manioc.RegisterScoped[IMyService, MyService](manioc.WithContainer(ctr)))

	// open scope
	scope, cleanup := manioc.OpenScope(manioc.WithParentScope(ctr))
	// resolve within the scope
	ret, err := manioc.Resolve[IMyService](manioc.WithScope(scope))
	assert.NotNil(ret)
	assert.Nil(err)

	// open child scope
	childScope, childCleanup := manioc.OpenScope(manioc.WithParentScope(scope))
	// resolve within the child scope
	childRet, childErr := manioc.Resolve[IMyService](manioc.WithScope(childScope))
	assert.NotNil(childRet)
	assert.Nil(childErr)

	// Since a scope does not inherit the instance cache,
	// resolution in different scopes would produce different instances.
	// We would like to make this behavior configurable.
	assert.NotSame(ret, childRet)

	// close scope
	cleanup()

	// Currently, when the parent scope is closed, the child scope will not be closed.
	// This behavior may be inappropriate if the cache is inherited.
	_, childErr = manioc.Resolve[IMyService](manioc.WithScope(childScope))
	assert.Nil(childErr)

	// close child scope
	childCleanup()
}
