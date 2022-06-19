package manioc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Activator_Errors(t *testing.T) {
	t.Run("field injection for non struct type", func(t *testing.T) {
		ctx := newDefaultContext()
		assert.Error(t, injectToFields(ctx, 42))
	})

	t.Run("activator error on NeverCache policy", func(t *testing.T) {
		assert := assert.New(t)

		ctr := NewContainer()

		assert.Nil(registerActivator[any](func(ctx resolveContext) (any, error) {
			return nil, fmt.Errorf("error")
		}, WithContainer(ctr)))
		assert.True(IsRegistered[any](WithContainer(ctr)))

		// resolution will fail since the registered activator is always returns error
		_, err := Resolve[any](WithScope(ctr))
		assert.Error(err)

		// also it will fail for resolve many
		_, err = Resolve[[]any](WithScope(ctr))
		assert.Error(err)
	})

	t.Run("activator error on ScopedCache policy", func(t *testing.T) {
		assert := assert.New(t)

		ctr := NewContainer()

		assert.Nil(registerActivator[any](func(ctx resolveContext) (any, error) {
			return nil, fmt.Errorf("error")
		}, WithContainer(ctr), WithCachePolicy(ScopedCache)))
		assert.True(IsRegistered[any](WithContainer(ctr)))

		// resolution will fail since the registered activator is always returns error
		_, err := Resolve[any](WithScope(ctr))
		assert.Error(err)
	})
}
