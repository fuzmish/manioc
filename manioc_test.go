package manioc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fieldInjection_for_non_struct(t *testing.T) {
	ctx := newDefaultContext()
	err := injectToFields(ctx, 42)
	assert.Error(t, err)
}

func Test_MalformedActivator(t *testing.T) {
	assert := assert.New(t)

	ctr := NewContainer()

	err := registerActivator[any](func(ctx resolveContext) (any, error) {
		return nil, fmt.Errorf("error")
	}, WithContainer(ctr), WithRegisterKey("T1"), WithCachePolicy(GlobalCache))
	assert.Nil(err)

	err = registerActivator[any](func(ctx resolveContext) (any, error) {
		return nil, fmt.Errorf("error")
	}, WithContainer(ctr), WithRegisterKey("T2"), WithCachePolicy(ScopedCache))
	assert.Nil(err)

	// resolve
	ret, err := Resolve[any](WithScope(ctr), WithResolveKey("T1"))
	assert.Nil(ret)
	assert.Error(err)

	ret, err = Resolve[any](WithScope(ctr), WithResolveKey("T2"))
	assert.Nil(ret)
	assert.Error(err)

	rets, err := Resolve[[]any](WithScope(ctr), WithResolveKey("T2"))
	assert.Nil(rets)
	assert.Error(err)
}
