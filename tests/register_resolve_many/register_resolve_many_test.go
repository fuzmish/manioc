package manioc_register_resolve_many_test

import (
	"reflect"
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

type MyService3 struct{}

func (s *MyService3) doSomething() {}

// utility functions
func getTypes[T any](data []T) []reflect.Type {
	if data == nil {
		return nil
	}
	ret := make([]reflect.Type, len(data))
	for i, elem := range data {
		ret[i] = reflect.TypeOf(elem)
	}
	return ret
}

func Test_RegisterAndResolveMany(t *testing.T) {
	t.Run("register one implementations", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register one implementations
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))

		// resolve all implementations
		services, err := manioc.ResolveMany[IMyService](manioc.WithScope(ctr))
		assert.Nil(err)
		assert.Len(services, 1)

		_, ok := services[0].((*MyService1))
		assert.True(ok)
	})

	t.Run("register two implementations", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register multiple implementations
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, MyService2](manioc.WithContainer(ctr)))

		// resolve all implementations
		services, err := manioc.ResolveMany[IMyService](manioc.WithScope(ctr))
		assert.Nil(err)
		assert.Len(services, 2)

		assert.ElementsMatch(
			getTypes([]IMyService{&MyService1{}, &MyService2{}}),
			getTypes(services),
		)
	})

	t.Run("ResolveMany[T] is an alias of Resolve[[]T]", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register multiple implementations
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, MyService2](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, MyService3](manioc.WithContainer(ctr)))

		// resolve all implementations
		services, err := manioc.Resolve[[]IMyService](manioc.WithScope(ctr))
		assert.Nil(err)
		assert.Len(services, 3)

		assert.ElementsMatch(
			getTypes([]IMyService{&MyService1{}, &MyService2{}, &MyService3{}}),
			getTypes(services),
		)
	})
}

func Test_ResolveMany_Errors(t *testing.T) {
	t.Run("if no registration found, returns error, not empty slice", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_, err := manioc.ResolveMany[IMyService](manioc.WithScope(ctr))
		assert.Error(err)
	})

	t.Run("Resolve will fail if multiple implementations are registered", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register multiple implementations
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, MyService2](manioc.WithContainer(ctr)))

		// resolve single will fail
		_, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
		assert.Error(err)
	})
}

func Test_MustResolveMany(t *testing.T) {
	t.Run("MustResolveMany causes panic for resolution errors", func(t *testing.T) {
		assert.Panics(t, func() {
			ctr := manioc.NewContainer()
			_ = manioc.MustResolveMany[IMyService](manioc.WithScope(ctr))
		})
	})

	t.Run("MustResolveMany returns the instance when resolution is successfully finished", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IMyService, MyService1](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, MyService2](manioc.WithContainer(ctr)))

		services := manioc.MustResolveMany[IMyService](manioc.WithScope(ctr))
		assert.Len(services, 2)
		assert.ElementsMatch(
			getTypes([]IMyService{&MyService1{}, &MyService2{}}),
			getTypes(services),
		)
	})
}
