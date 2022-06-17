package manioc_test

import (
	"testing"

	"github.com/fuzmish/manioc"
	"github.com/stretchr/testify/assert"
)

type IFooService interface {
	DoFoo()
}

type FooService struct {
	Foo int
}

func NewFooService() *FooService {
	return &FooService{Foo: 42}
}

func (s *FooService) DoFoo() {
}

type IBarService interface {
	DoBar()
}

type BarService struct {
	skfoo IFooService
}

func (s *BarService) DoBar() {
}

func NewBarService(foo IFooService) *BarService {
	return &BarService{skfoo: foo}
}

type BarService2 struct{}

func (s *BarService2) DoBar() {
}

type IBazService interface {
	DoBaz()
}

type BazService struct {
	Skbar  IBarService `manioc:"inject"`
	Skbar2 IBarService `manioc:"inject,key=hello"`
}

func (s *BazService) DoBaz() {
	s.Skbar.DoBar()
	s.Skbar2.DoBar()
}

type IFooBarService interface {
	DoFooBar() int
}

type FooBarService struct {
	JsBar []IBarService `manioc:"inject"`
}

func NewFooBarService() *FooBarService {
	return &FooBarService{}
}

func (s *FooBarService) DoFooBar() int {
	return len(s.JsBar)
}

type IFooBazService interface {
	DoFooBaz() int
}

type FooBazService struct {
	xsBar []IBarService
}

func (s *FooBazService) DoFooBaz() int {
	return len(s.xsBar)
}

func NewFooBazService(bar []IBarService) *FooBazService {
	return &FooBazService{xsBar: bar}
}

type IBarBazService interface {
	DoBarBaz()
}

type BarBazService struct {
	foo IFooService `manioc:"inject"`
}

func (s *BarBazService) DoBarBaz() {
	s.foo.DoFoo()
}

func Test_Fail(t *testing.T) {
	t.Run("Invalid CachePolicy", func(t *testing.T) {
		assert := assert.New(t)
		assert.Panics(func() {
			ctr := manioc.NewContainer()
			err := manioc.Register[IFooService, FooService](manioc.WithContainer(ctr), manioc.WithCachePolicy(42))
			assert.Nil(err)
		})
	})

	t.Run("resolve implementation before registration should fail", func(t *testing.T) {
		assert := assert.New(t)
		ctr := manioc.NewContainer()
		ret, err := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		assert.Nil(ret)
		assert.Error(err)
	})

	t.Run("resolve implementation with field injections before registration should fail", func(t *testing.T) {
		assert := assert.New(t)
		ctr := manioc.NewContainer()
		err := manioc.Register[IBazService, BazService](manioc.WithContainer(ctr))
		assert.Nil(err)
		ret2, err := manioc.Resolve[IBazService](manioc.WithScope(ctr))
		assert.Nil(ret2)
		assert.Error(err)
	})

	t.Run("resolve constructor before registration should fail", func(t *testing.T) {
		assert := assert.New(t)
		ctr := manioc.NewContainer()
		err := manioc.RegisterConstructor[IFooBazService](NewFooBazService, manioc.WithContainer(ctr))
		assert.Nil(err)
		ret4, err := manioc.Resolve[IFooBazService](manioc.WithScope(ctr))
		assert.Nil(ret4)
		assert.Error(err)
	})

	t.Run("resolve constructor with field injections before registration should fail", func(t *testing.T) {
		assert := assert.New(t)
		ctr := manioc.NewContainer()
		err := manioc.RegisterConstructor[IFooBarService](NewFooBarService, manioc.WithContainer(ctr))
		assert.Nil(err)
		ret4, err := manioc.Resolve[IFooBarService](manioc.WithScope(ctr))
		assert.Nil(ret4)
		assert.Error(err)

	})
}

func Test_Fail2(t *testing.T) {
	assert := assert.New(t)

	ctr := manioc.NewContainer()

	err := manioc.Register[IBarService, BarService](manioc.WithContainer(ctr))
	assert.Nil(err)
	err = manioc.Register[IBarService, BarService2](manioc.WithContainer(ctr))
	assert.Nil(err)

	ret, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
	assert.Nil(ret)
	assert.Error(err)
}

func Test_InvalidRegistration(t *testing.T) {
	t.Run("invalid implementation type", func(t *testing.T) {
		assert.Panics(t, func() {
			ctr := manioc.NewContainer()
			err := manioc.Register[IFooService, *struct{}](manioc.WithContainer(ctr))
			assert.Error(t, err)
		})
	})

	t.Run("invalid constructor type", func(t *testing.T) {
		assert.Panics(t, func() {
			ctr := manioc.NewContainer()
			err := manioc.RegisterConstructor[IFooService](func() *struct{} { return nil }, manioc.WithContainer(ctr))
			assert.Error(t, err)
		})
	})
}

func Test_Container(t *testing.T) {
	ctr := manioc.NewContainer()
	err := manioc.Register[IFooService, FooService](
		manioc.WithContainer(ctr),
		manioc.WithCachePolicy(manioc.GlobalCache),
		manioc.WithRegisterKey(42))
	if err != nil {
		t.FailNow()
	}
	ret, err := manioc.Resolve[IFooService](
		manioc.WithScope(ctr),
		manioc.WithResolveKey(42))
	if err != nil {
		t.FailNow()
	}
	ret.DoFoo()
}

func Test_Transient(t *testing.T) {
	t.Run("default cache policy is NeverCache", func(t *testing.T) {
		ctr := manioc.NewContainer()
		err := manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))
		assert.Nil(t, err)
		r1, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		r2, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		assert.NotSame(t, r1, r2)
	})

	// same test but using explicit helper API
	t.Run("traisient", func(t *testing.T) {
		ctr := manioc.NewContainer()
		err := manioc.RegisterTransient[IFooService, FooService](manioc.WithContainer(ctr))
		assert.Nil(t, err)
		r1, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		r2, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		assert.NotSame(t, r1, r2)
	})

	// for constructor
	t.Run("traisient constructor", func(t *testing.T) {
		ctr := manioc.NewContainer()
		err := manioc.RegisterTransientConstructor[IFooService](NewFooService, manioc.WithContainer(ctr))
		assert.Nil(t, err)
		r1, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		r2, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		assert.NotSame(t, r1, r2)
	})
}

func Test_Scope(t *testing.T) {
	t.Run("using ScopedCache policy", func(t *testing.T) {
		ctr := manioc.NewContainer()
		err := manioc.Register[IFooService, FooService](
			manioc.WithCachePolicy(manioc.ScopedCache),
			manioc.WithContainer(ctr),
		)
		assert.Nil(t, err)
		scope, cleanup := manioc.OpenScope(manioc.WithParentScope(ctr))
		r1, _ := manioc.Resolve[IFooService](manioc.WithScope(scope))
		r2, _ := manioc.Resolve[IFooService](manioc.WithScope(scope))
		assert.Same(t, r1, r2)
		cleanup()

		// after cleanup
		_, err = manioc.Resolve[IFooService](manioc.WithScope(scope))
		assert.Error(t, err)
	})

	// same test but using helper API
	t.Run("using ScopedCache policy", func(t *testing.T) {
		ctr := manioc.NewContainer()
		err := manioc.RegisterScoped[IFooService, FooService](manioc.WithContainer(ctr))
		assert.Nil(t, err)
		scope, cleanup := manioc.OpenScope(manioc.WithParentScope(ctr))
		r1, _ := manioc.Resolve[IFooService](manioc.WithScope(scope))
		r2, _ := manioc.Resolve[IFooService](manioc.WithScope(scope))
		assert.Same(t, r1, r2)
		cleanup()

		// after cleanup
		_, err = manioc.Resolve[IFooService](manioc.WithScope(scope))
		assert.Error(t, err)
	})

	// for constructor
	t.Run("scoped constructor", func(t *testing.T) {
		ctr := manioc.NewContainer()
		err := manioc.RegisterScopedConstructor[IFooService](NewFooService, manioc.WithContainer(ctr))
		assert.Nil(t, err)
		scope, cleanup := manioc.OpenScope(manioc.WithParentScope(ctr))
		r1, _ := manioc.Resolve[IFooService](manioc.WithScope(scope))
		r2, _ := manioc.Resolve[IFooService](manioc.WithScope(scope))
		assert.Same(t, r1, r2)
		cleanup()

		// after cleanup
		_, err = manioc.Resolve[IFooService](manioc.WithScope(scope))
		assert.Error(t, err)
	})
}

func Test_Singleton(t *testing.T) {
	t.Run("using GlobalCache policy", func(t *testing.T) {
		ctr := manioc.NewContainer()
		err := manioc.Register[IFooService, FooService](
			manioc.WithContainer(ctr),
			manioc.WithCachePolicy(manioc.GlobalCache),
		)
		assert.Nil(t, err)
		r1, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		r2, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		assert.Same(t, r1, r2)
	})

	t.Run("using GlobalCache policy", func(t *testing.T) {
		ctr := manioc.NewContainer()
		err := manioc.RegisterSingleton[IFooService, FooService](manioc.WithContainer(ctr))
		assert.Nil(t, err)
		r1, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		r2, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		assert.Same(t, r1, r2)
	})

	t.Run("using GlobalCache policy", func(t *testing.T) {
		ctr := manioc.NewContainer()
		err := manioc.RegisterSingletonConstructor[IFooService](NewFooService, manioc.WithContainer(ctr))
		assert.Nil(t, err)
		r1, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		r2, _ := manioc.Resolve[IFooService](manioc.WithScope(ctr))
		assert.Same(t, r1, r2)
	})
}

func Test_RegisterConstructor(t *testing.T) {
	ctr := manioc.NewContainer()

	err := manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))
	assert.Nil(t, err)
	err = manioc.RegisterConstructor[IBarService](NewBarService, manioc.WithContainer(ctr))
	assert.Nil(t, err)

	// resolve
	ret, err := manioc.Resolve[IBarService](manioc.WithScope(ctr))
	if err != nil {
		t.FailNow()
	}
	ret.DoBar()
}

func Test_FieldInjection(t *testing.T) {
	ctr := manioc.NewContainer()

	err := manioc.RegisterInstance[IFooService](&FooService{}, manioc.WithContainer(ctr))
	assert.Nil(t, err)
	err = manioc.RegisterConstructor[IBarService](NewBarService, manioc.WithContainer(ctr))
	assert.Nil(t, err)
	err = manioc.RegisterConstructor[IBarService](NewBarService, manioc.WithContainer(ctr), manioc.WithRegisterKey("hello"))
	assert.Nil(t, err)
	err = manioc.Register[IBazService, BazService](manioc.WithContainer(ctr))
	assert.Nil(t, err)

	// resolve
	ret, err := manioc.Resolve[IBazService](manioc.WithScope(ctr))
	if err != nil {
		t.FailNow()
	}
	ret.DoBaz()
}

func Test_ResolveMany(t *testing.T) {
	ctr := manioc.NewContainer()

	err := manioc.Register[IBarService, BarService](manioc.WithContainer(ctr))
	assert.Nil(t, err)
	err = manioc.Register[IBarService, BarService](manioc.WithContainer(ctr), manioc.WithRegisterKey("bar"))
	assert.Nil(t, err)
	err = manioc.Register[IBarService, BarService2](manioc.WithContainer(ctr), manioc.WithRegisterKey("bar"))
	assert.Nil(t, err)

	ret, err := manioc.ResolveMany[IBarService](manioc.WithScope(ctr))
	if err != nil {
		t.FailNow()
	}
	if len(ret) != 1 {
		t.FailNow()
	}
	retbar, err := manioc.ResolveMany[IBarService](manioc.WithScope(ctr), manioc.WithResolveKey("bar"))
	if err != nil {
		t.FailNow()
	}
	if len(retbar) != 2 {
		t.FailNow()
	}
	retbardirect, err := manioc.Resolve[[]IBarService](manioc.WithScope(ctr), manioc.WithResolveKey("bar"))
	if err != nil {
		t.FailNow()
	}
	if len(retbardirect) != 2 {
		t.FailNow()
	}
}

func Test_ResolveManyViaFieldInjection(t *testing.T) {
	ctr := manioc.NewContainer()

	err := manioc.Register[IBarService, BarService](manioc.WithContainer(ctr))
	assert.Nil(t, err)
	err = manioc.Register[IBarService, BarService2](manioc.WithContainer(ctr))
	assert.Nil(t, err)
	err = manioc.Register[IFooBarService, FooBarService](manioc.WithContainer(ctr))
	assert.Nil(t, err)

	ret, err := manioc.Resolve[IFooBarService](manioc.WithScope(ctr))
	if err != nil {
		t.FailNow()
	}
	if ret.DoFooBar() != 2 {
		t.FailNow()
	}
}

func Test_ResolveManyViaConstructorInjection(t *testing.T) {
	ctr := manioc.NewContainer()

	err := manioc.Register[IBarService, BarService](manioc.WithContainer(ctr))
	assert.Nil(t, err)
	err = manioc.Register[IBarService, BarService2](manioc.WithContainer(ctr))
	assert.Nil(t, err)
	err = manioc.RegisterConstructor[IFooBazService](NewFooBazService, manioc.WithContainer(ctr))
	assert.Nil(t, err)

	ret, err := manioc.Resolve[IFooBazService](manioc.WithScope(ctr))
	if err != nil {
		t.FailNow()
	}
	if ret.DoFooBaz() != 2 {
		t.FailNow()
	}
}

func Test_unexported_field_inject(t *testing.T) {
	assert := assert.New(t)

	ctr := manioc.NewContainer()

	err := manioc.Register[IFooService, FooService](manioc.WithContainer(ctr))
	assert.Nil(err)
	err = manioc.Register[IBarBazService, BarBazService](manioc.WithContainer(ctr))
	assert.Nil(err)

	ret, err := manioc.Resolve[IBarBazService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)
	ret.DoBarBaz()
}

func Test_IsRegistered(t *testing.T) {
	assert := assert.New(t)

	ctr := manioc.NewContainer()

	// initial state; the container have no registration
	assert.False(manioc.IsRegistered[IFooService](manioc.WithContainer(ctr)))
	assert.False(manioc.IsRegistered[IFooService](manioc.WithContainer(ctr), manioc.WithRegisterKey("foo")))

	// register
	assert.Nil(manioc.Register[IFooService, FooService](manioc.WithContainer(ctr)))

	// check registration
	assert.True(manioc.IsRegistered[IFooService](manioc.WithContainer(ctr)))
	assert.False(manioc.IsRegistered[IFooService](manioc.WithContainer(ctr), manioc.WithRegisterKey("foo")))

	// register with key
	assert.Nil(manioc.Register[IFooService, FooService](manioc.WithContainer(ctr), manioc.WithRegisterKey("foo")))

	// check registration
	assert.True(manioc.IsRegistered[IFooService](manioc.WithContainer(ctr)))
	assert.True(manioc.IsRegistered[IFooService](manioc.WithContainer(ctr), manioc.WithRegisterKey("foo")))

	// resolution will success
	_, err := manioc.Resolve[IFooService](manioc.WithScope(ctr))
	assert.Nil(err)
	_, err = manioc.Resolve[IFooService](manioc.WithScope(ctr), manioc.WithResolveKey("foo"))
	assert.Nil(err)

	// unregister not registered interface
	assert.False(manioc.Unregister[IBarService](manioc.WithContainer(ctr)))
	assert.False(manioc.Unregister[IBarService](manioc.WithContainer(ctr), manioc.WithRegisterKey("bar")))

	// unregister nil key registration
	assert.True(manioc.Unregister[IFooService](manioc.WithContainer(ctr)))
	// already unregistered
	assert.False(manioc.Unregister[IFooService](manioc.WithContainer(ctr)))

	// check registration
	assert.False(manioc.IsRegistered[IFooService](manioc.WithContainer(ctr)))
	assert.True(manioc.IsRegistered[IFooService](manioc.WithContainer(ctr), manioc.WithRegisterKey("foo")))

	// unregister keyed registration
	assert.True(manioc.Unregister[IFooService](manioc.WithContainer(ctr), manioc.WithRegisterKey("foo")))
	// already unregistered
	assert.False(manioc.Unregister[IFooService](manioc.WithContainer(ctr), manioc.WithRegisterKey("foo")))

	// check registration
	assert.False(manioc.IsRegistered[IFooService](manioc.WithContainer(ctr)))
	assert.False(manioc.IsRegistered[IFooService](manioc.WithContainer(ctr), manioc.WithRegisterKey("foo")))

	// resolution will fail
	_, err = manioc.Resolve[IFooService](manioc.WithScope(ctr))
	assert.Error(err)
	_, err = manioc.Resolve[IFooService](manioc.WithScope(ctr), manioc.WithResolveKey("foo"))
	assert.Error(err)
}
