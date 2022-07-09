package manioc_empty_interface_test

import (
	"testing"

	"github.com/fuzmish/manioc"
	"github.com/stretchr/testify/assert"
)

type IMyService interface{}

type MyService struct{}

type MyService2 struct {
	Value int
}

func (s *MyService2) DoSomething() {}

func Test_RegistrationForEmptyInterface(t *testing.T) {
	t.Run("any accepts any implementation types", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[any, bool](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[any, int](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[any, struct{}](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[any, MyService](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[any, MyService2](manioc.WithContainer(ctr)))
	})

	t.Run("empty interface accepts any implementation types", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IMyService, bool](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, int](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, struct{}](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, MyService](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, MyService2](manioc.WithContainer(ctr)))
	})
}

func Test_ResolveWithEmptyInterface(t *testing.T) {
	t.Run("an empty interface is still distinguished by its type", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[any, MyService](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IMyService, MyService2](manioc.WithContainer(ctr)))

		ret, err := manioc.Resolve[any](manioc.WithScope(ctr))
		assert.Nil(err)
		assert.IsType(&MyService{}, ret)
		ret, err = manioc.Resolve[IMyService](manioc.WithScope(ctr))
		assert.Nil(err)
		assert.IsType(&MyService2{}, ret)
	})
}
