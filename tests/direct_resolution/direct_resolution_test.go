package manioc_direct_resolution_test

import (
	"errors"
	"testing"

	"github.com/fuzmish/manioc"
	"github.com/stretchr/testify/assert"
)

type IFooService interface {
	doFoo()
}

type FooService struct{}

func (s *FooService) doFoo() {}

type BarService struct {
	foo IFooService `manioc:"inject"`
}

type BazService struct {
	foo1 IFooService
	foo2 IFooService `manioc:"inject"`
}

func NewBazService(foo IFooService) *BazService {
	return &BazService{foo1: foo}
}

func NewBazServiceWithError(foo IFooService) (*BazService, error) {
	return &BazService{foo1: foo}, nil
}

func NewBazServiceRaiseError(foo IFooService) (*BazService, error) {
	return nil, errors.New("error")
}

func Test_ResolveInstance(t *testing.T) {
	t.Run("basic case", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))

		ret, err := manioc.ResolveInstance(&BarService{}, manioc.WithScope(ctr))
		assert.Nil(err)
		assert.NotNil(ret)
		assert.NotNil(ret.foo)
	})

	t.Run("dependency is missing", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		_, err := manioc.ResolveInstance(&BarService{}, manioc.WithScope(ctr))
		assert.Error(err)
	})

	t.Run("cannot resolve nil value", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))

		_, err := manioc.ResolveInstance[*BarService](nil, manioc.WithScope(ctr))
		assert.Error(err)
	})
}

func Test_MustResolveInstance(t *testing.T) {
	t.Run("basic case", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))

		ret := manioc.MustResolveInstance(&BarService{}, manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.NotNil(ret.foo)
	})

	t.Run("dependency is missing", func(t *testing.T) {
		assert.Panics(t, func() {
			ctr := manioc.NewContainer()
			_ = manioc.MustResolveInstance(&BarService{}, manioc.WithScope(ctr))
		})
	})

	t.Run("cannot resolve nil value", func(t *testing.T) {
		assert.Panics(t, func() {
			ctr := manioc.NewContainer()
			_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))
			_ = manioc.MustResolveInstance[*BarService](nil, manioc.WithScope(ctr))
		})
	})
}

func Test_ResolveFunction(t *testing.T) {
	t.Run("basic case", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))

		ret, err := manioc.ResolveFunction[*BazService](NewBazService, manioc.WithScope(ctr))
		assert.Nil(err)
		assert.NotNil(ret)
		assert.NotNil(ret.foo1)
		assert.NotNil(ret.foo2)
	})

	t.Run("dependency is missing", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		_, err := manioc.ResolveFunction[*BazService](NewBazService, manioc.WithScope(ctr))
		assert.Error(err)
	})

	t.Run("function with error return value", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))

		ret, err := manioc.ResolveFunction[*BazService](NewBazServiceWithError, manioc.WithScope(ctr))
		assert.Nil(err)
		assert.NotNil(ret)
		assert.NotNil(ret.foo1)
		assert.NotNil(ret.foo2)
	})

	t.Run("error propagation from function", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))

		_, err := manioc.ResolveFunction[*BazService](NewBazServiceRaiseError, manioc.WithScope(ctr))
		assert.Error(err)
	})

	t.Run("cannot resolve nil value", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))

		_, err := manioc.ResolveFunction[*BazService, func() *BazService](nil, manioc.WithScope(ctr))
		assert.Error(err)
	})
}

func Test_MustResolveFunction(t *testing.T) {
	t.Run("basic case", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))

		ret := manioc.MustResolveFunction[*BazService](NewBazService, manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.NotNil(ret.foo1)
		assert.NotNil(ret.foo2)
	})

	t.Run("dependency is missing", func(t *testing.T) {
		assert.Panics(t, func() {
			ctr := manioc.NewContainer()
			_ = manioc.MustResolveFunction[*BazService](NewBazService, manioc.WithScope(ctr))
		})
	})

	t.Run("function with error return value", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))

		ret := manioc.MustResolveFunction[*BazService](NewBazServiceWithError, manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.NotNil(ret.foo1)
		assert.NotNil(ret.foo2)
	})

	t.Run("error propagation from function", func(t *testing.T) {
		assert.Panics(t, func() {
			ctr := manioc.NewContainer()
			_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))
			_ = manioc.MustResolveFunction[*BazService](NewBazServiceRaiseError, manioc.WithScope(ctr))
		})
	})

	t.Run("cannot resolve nil value", func(t *testing.T) {
		assert.Panics(t, func() {
			ctr := manioc.NewContainer()
			_ = manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))
			_ = manioc.MustResolveFunction[*BazService, func() *BazService](nil, manioc.WithScope(ctr))
		})
	})
}
