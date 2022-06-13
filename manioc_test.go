package manioc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MalformedActivator(t *testing.T) {
	assert := assert.New(t)

	ctr := NewContainer()

	err := registerActivator[any](func(ctx resolveContext) (any, error) {
		return nil, fmt.Errorf("ERROR!")
	}, WithContainer(ctr), WithRegisterKey("T1"), WithCachePolicy(GlobalCache))
	assert.Nil(err)

	err = registerActivator[any](func(ctx resolveContext) (any, error) {
		return nil, fmt.Errorf("ERROR!")
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
