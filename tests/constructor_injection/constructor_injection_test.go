package manioc_constructor_injection_test

import (
	"testing"

	"github.com/fuzmish/manioc"
	"github.com/stretchr/testify/assert"
)

type IFooService interface {
	doFoo()
}

type FooService1 struct{}

func (s *FooService1) doFoo() {}

type FooService2 struct{}

func (s *FooService2) doFoo() {}

type IBarService interface {
	doBar()
}

type BarService struct {
	foo []IFooService
}

func (s *BarService) doBar() {}

func NewBarService(foo IFooService) *BarService {
	return &BarService{foo: []IFooService{foo}}
}

func NewBarServiceWithManyFoo(foo []IFooService) *BarService {
	return &BarService{foo: foo}
}

type IBazService interface {
	doBaz()
}

type BazService struct {
	foo IFooService
	bar IBarService
}

func (s *BazService) doBaz() {}

func NewBazService(foo IFooService, bar IBarService) *BazService {
	return &BazService{foo: foo, bar: bar}
}

func Test_ConstructorInjection(t *testing.T) {
	t.Run("resolve with constructor injection before registration", func(t *testing.T) {
		assert := assert.New(t)

		// register
		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterConstructor[IBarService](NewBarService, manioc.WithContainer(ctr)))

		// resolution will fail since the dependency IFooService is not registered
		_, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
		assert.Error(err)
	})

	t.Run("basic case", func(t *testing.T) {
		assert := assert.New(t)

		// register
		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
		assert.Nil(manioc.RegisterConstructor[IBarService](NewBarService, manioc.WithContainer(ctr)))

		// resolve
		ret, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Nil(err)

		// check injected instance
		bar, ok := ret.(*BarService)
		assert.True(ok)
		assert.Len(bar.foo, 1)
		_, ok = bar.foo[0].(*FooService1)
		assert.True(ok)
	})
}

func Test_ConstructorInjection_ResolveMany(t *testing.T) {
	assert := assert.New(t)

	// register
	ctr := manioc.NewContainer()
	assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
	assert.Nil(manioc.Register[IFooService, FooService2](manioc.WithContainer(ctr)))
	assert.Nil(manioc.RegisterConstructor[IBarService](NewBarServiceWithManyFoo, manioc.WithContainer(ctr)))

	// resolve
	ret, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)

	// check injected instances
	bar, ok := ret.(*BarService)
	assert.True(ok)
	assert.Len(bar.foo, 2)
}

func Test_ConstructorInjection_MultipleArgs(t *testing.T) {
	assert := assert.New(t)

	// register
	ctr := manioc.NewContainer()
	assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
	assert.Nil(manioc.Register[IBarService, BarService](manioc.WithContainer(ctr)))
	assert.Nil(manioc.RegisterConstructor[IBazService](NewBazService, manioc.WithContainer(ctr)))

	// resolve
	ret, err := manioc.Resolve[IBazService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)

	// check injected instances
	baz, ok := ret.(*BazService)
	assert.True(ok)
	assert.NotNil(baz.foo)
	assert.NotNil(baz.bar)
}

func Test_ConstructorInjection_Recursive(t *testing.T) {
	assert := assert.New(t)

	// register
	ctr := manioc.NewContainer()
	assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
	assert.Nil(manioc.RegisterConstructor[IBarService](NewBarService, manioc.WithContainer(ctr)))
	assert.Nil(manioc.RegisterConstructor[IBazService](NewBazService, manioc.WithContainer(ctr)))

	// resolve
	ret, err := manioc.Resolve[IBazService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)

	// check injected instances
	baz, ok := ret.(*BazService)
	assert.True(ok)
	assert.NotNil(baz.foo)
	assert.NotNil(baz.bar)
	bar, ok := baz.bar.(*BarService)
	assert.True(ok)
	assert.Len(bar.foo, 1)
}
