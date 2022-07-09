package manioc_field_injection_test

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

type BarServiceSimple struct {
	Foo IFooService `manioc:"inject"`
}

func (s *BarServiceSimple) doBar() {}

func NewBarServiceSimple() *BarServiceSimple {
	return &BarServiceSimple{}
}

type BarServiceUnexported struct {
	foo IFooService `manioc:"inject"`
}

func (s *BarServiceUnexported) doBar() {}

type BarServiceInvalidTag struct {
	//nolint:unused
	foo IFooService `manioc:"iject"` // typo
}

func (s *BarServiceInvalidTag) doBar() {}

type BarServiceWithManyFoo struct {
	foo []IFooService `manioc:"inject"`
}

func (s *BarServiceWithManyFoo) doBar() {}

type BarServiceWithResolveKey struct {
	foo        IFooService `manioc:"inject"`
	fooAnother IFooService `manioc:"inject,key=another"`
}

func (s *BarServiceWithResolveKey) doBar() {}

type IBazService interface {
	doBaz()
}

type BazService struct {
	bar IBarService `manioc:"inject"`
}

func (s *BazService) doBaz() {}

func Test_FieldInjection(t *testing.T) {
	t.Run("resolution with field injection before registration", func(t *testing.T) {
		assert := assert.New(t)

		// register
		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IBarService, BarServiceSimple](manioc.WithContainer(ctr)))

		// resolve
		_, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
		assert.Error(err)
	})

	t.Run("basic case", func(t *testing.T) {
		assert := assert.New(t)

		// register
		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
		assert.Nil(manioc.Register[IBarService, BarServiceSimple](manioc.WithContainer(ctr)))

		// resolve
		ret, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Nil(err)

		bar, ok := ret.(*BarServiceSimple)
		assert.True(ok)
		assert.NotNil(bar.Foo)
	})
}

func Test_FieldInjection_WithConstructor(t *testing.T) {
	t.Run("resolution with field injection before registration", func(t *testing.T) {
		assert := assert.New(t)

		// register
		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterConstructor[IBarService](NewBarServiceSimple, manioc.WithContainer(ctr)))

		// resolve
		_, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
		assert.Error(err)
	})

	t.Run("basic case", func(t *testing.T) {
		assert := assert.New(t)

		// register
		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
		assert.Nil(manioc.RegisterConstructor[IBarService](NewBarServiceSimple, manioc.WithContainer(ctr)))

		// resolve
		ret, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Nil(err)

		bar, ok := ret.(*BarServiceSimple)
		assert.True(ok)
		assert.NotNil(bar.Foo)
	})
}

func Test_FieldInjection_UnexportedField(t *testing.T) {
	assert := assert.New(t)

	// register
	ctr := manioc.NewContainer()
	assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
	assert.Nil(manioc.Register[IBarService, BarServiceUnexported](manioc.WithContainer(ctr)))

	// resolve
	ret, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)

	bar, ok := ret.(*BarServiceUnexported)
	assert.True(ok)
	assert.NotNil(bar.foo)
}

func Test_FieldInjection_InvalidTag(t *testing.T) {
	assert := assert.New(t)

	// register
	ctr := manioc.NewContainer()
	assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
	assert.Nil(manioc.Register[IBarService, BarServiceInvalidTag](manioc.WithContainer(ctr)))

	// resolve
	_, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
	assert.Error(err)
}

func Test_FieldInjection_ResolveMany(t *testing.T) {
	assert := assert.New(t)

	// register
	ctr := manioc.NewContainer()
	assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
	assert.Nil(manioc.Register[IFooService, FooService2](manioc.WithContainer(ctr)))
	assert.Nil(manioc.Register[IBarService, BarServiceWithManyFoo](manioc.WithContainer(ctr)))

	// resolve
	ret, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)

	// check injected instances
	bar, ok := ret.(*BarServiceWithManyFoo)
	assert.True(ok)
	assert.Len(bar.foo, 2)
}

func Test_FieldInjection_WithKey(t *testing.T) {
	assert := assert.New(t)

	// register
	ctr := manioc.NewContainer()
	assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
	assert.Nil(manioc.Register[IFooService, FooService2](
		manioc.WithContainer(ctr),
		manioc.WithRegisterKey("another"),
	))
	assert.Nil(manioc.Register[IBarService, BarServiceWithResolveKey](manioc.WithContainer(ctr)))

	// resolve
	ret, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)

	bar, ok := ret.(*BarServiceWithResolveKey)
	assert.True(ok)
	_, ok = bar.foo.(*FooService1)
	assert.True(ok)
	_, ok = bar.fooAnother.(*FooService2)
	assert.True(ok)
}

func Test_FieldInjection_RegisteredInstance(t *testing.T) {
	t.Run("resolution with field injection before registration", func(t *testing.T) {
		assert := assert.New(t)

		// register
		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterInstance[IBarService](&BarServiceSimple{}, manioc.WithContainer(ctr)))

		// resolve
		_, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
		assert.Error(err)
	})

	t.Run("basic case", func(t *testing.T) {
		assert := assert.New(t)

		// register
		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
		assert.Nil(manioc.RegisterInstance[IBarService](&BarServiceSimple{}, manioc.WithContainer(ctr)))

		// resolve
		ret, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Nil(err)

		bar, ok := ret.(*BarServiceSimple)
		assert.True(ok)
		assert.NotNil(bar.Foo)
	})
}

func Test_FieldInjection_Recursive(t *testing.T) {
	assert := assert.New(t)

	// register
	ctr := manioc.NewContainer()
	assert.Nil(manioc.Register[IFooService, FooService1](manioc.WithContainer(ctr)))
	assert.Nil(manioc.Register[IBarService, BarServiceSimple](manioc.WithContainer(ctr)))
	assert.Nil(manioc.Register[IBazService, BazService](manioc.WithContainer(ctr)))

	// resolve
	ret, err := manioc.Resolve[IBazService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)

	baz, ok := ret.(*BazService)
	assert.True(ok)
	bar, ok := baz.bar.(*BarServiceSimple)
	assert.True(ok)
	_, ok = bar.Foo.(*FooService1)
	assert.True(ok)
}
