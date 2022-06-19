package manioc_query_registration_test

import (
	"testing"

	"github.com/fuzmish/manioc"
	"github.com/stretchr/testify/assert"
)

type IMyService interface {
	doSomething()
}

// MyServiceN implements IMyService
type MyService1 struct{}

func (s *MyService1) doSomething() {}

type MyService2 struct{}

func (s *MyService2) doSomething() {}

func Test_IsRegistered(t *testing.T) {
	t.Run("just after the container was created, no implementation is regitered", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.False(manioc.IsRegistered[any](manioc.WithContainer(ctr)))
		assert.False(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))
	})

	t.Run("query registry", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register service
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))

		// verify that it was registered correctly
		assert.True(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))
	})

	t.Run("multiple registration", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register service
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, MyService2](manioc.WithContainer(ctr)))

		// verify that it was registered correctly
		assert.True(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))
	})
}

func Test_IsRegistered_WithKey(t *testing.T) {
	t.Run("query registry with key", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register service
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))

		// verify that it was registered correctly
		assert.True(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))

		// since it is currently registered for anonymous service key,
		// it appears unregistered when the key is specified in a query.
		assert.False(manioc.IsRegistered[IMyService](
			manioc.WithContainer(ctr),
			manioc.WithRegisterKey("mykey")))
	})

	t.Run("query registry with registered key", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register service with key
		key := "mykey"
		assert.Nil(manioc.Register[IMyService, MyService1](
			manioc.WithContainer(ctr),
			manioc.WithRegisterKey(key)))

		// verify that it was registered correctly
		assert.True(manioc.IsRegistered[IMyService](
			manioc.WithContainer(ctr),
			manioc.WithRegisterKey("mykey")))

		// since it is currently registered with key,
		// it appears unregistered for anonymous key.
		assert.False(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))
	})
}

func Test_Unregister(t *testing.T) {
	t.Run("if no registration found, Unregister returns false", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.False(manioc.Unregister[IMyService](manioc.WithContainer(ctr)))
	})

	t.Run("basic usage", func(t *testing.T) {
		assert := assert.New(t)

		// init container
		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))
		assert.True(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))

		// unregister
		assert.True(manioc.Unregister[IMyService](manioc.WithContainer(ctr)))

		// verify
		assert.False(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))
	})

	t.Run("when multiple implementation is registered", func(t *testing.T) {
		assert := assert.New(t)

		// init container
		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, MyService2](manioc.WithContainer(ctr)))
		assert.True(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))

		// unregister; this operation removes all registration for IMyService
		assert.True(manioc.Unregister[IMyService](manioc.WithContainer(ctr)))

		// verify
		assert.False(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))
	})
}

func Test_Unregister_WithKey(t *testing.T) {
	t.Run("unregister with key", func(t *testing.T) {
		assert := assert.New(t)

		// init container
		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))
		assert.True(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))

		// since it is currently registered for anonymous service key,
		// unregistration with key will fail.
		assert.False(manioc.Unregister[IMyService](manioc.WithContainer(ctr), manioc.WithRegisterKey("mykey")))

		// the registration status does not change
		assert.True(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))

		// unregister
		assert.True(manioc.Unregister[IMyService](manioc.WithContainer(ctr)))

		// verify
		assert.False(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))
	})

	t.Run("unregister with registered key", func(t *testing.T) {
		assert := assert.New(t)

		// init container
		ctr := manioc.NewContainer()
		key := "mykey"
		assert.Nil(manioc.Register[IMyService, MyService1](
			manioc.WithContainer(ctr),
			manioc.WithRegisterKey(key),
		))
		assert.True(manioc.IsRegistered[IMyService](
			manioc.WithContainer(ctr),
			manioc.WithRegisterKey(key),
		))

		// since it is currently registered with key,
		// unregistration with anonymous key will fail.
		assert.False(manioc.Unregister[IMyService](manioc.WithContainer(ctr)))

		// the registration status does not change
		assert.True(manioc.IsRegistered[IMyService](
			manioc.WithContainer(ctr),
			manioc.WithRegisterKey(key),
		))

		// unregister
		assert.True(manioc.Unregister[IMyService](
			manioc.WithContainer(ctr),
			manioc.WithRegisterKey(key),
		))

		// verify
		assert.False(manioc.IsRegistered[IMyService](manioc.WithContainer(ctr)))
	})
}
