package manioc_register_and_resolve_test

import (
	"errors"
	"testing"

	"github.com/fuzmish/manioc"
	"github.com/stretchr/testify/assert"
)

type IMyService interface {
	doSomething()
}

// MyService implements IMyService
type MyService struct{}

func (s *MyService) doSomething() {}

func NewMyService() *MyService {
	return &MyService{}
}

func NewMyServiceWithError() (*MyService, error) {
	return &MyService{}, nil
}

func NewMyServiceRaiseError() (*MyService, error) {
	return nil, errors.New("error")
}

func Test_RegisterAndResolve(t *testing.T) {
	assert := assert.New(t)

	ctr := manioc.NewContainer()

	// register
	assert.Nil(manioc.Register[IMyService, MyService](manioc.WithContainer(ctr)))

	// resolve
	ret, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)

	// verify that the value of ret is the instance of MyService
	_, ok := ret.(*MyService)
	assert.True(ok)
}

func Test_RegisterConstructorAndResolve(t *testing.T) {
	t.Run("constructor without error", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register
		assert.Nil(manioc.RegisterConstructor[IMyService](NewMyService, manioc.WithContainer(ctr)))

		// resolve
		ret, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Nil(err)

		// verify that the value of ret is the instance of MyService
		_, ok := ret.(*MyService)
		assert.True(ok)
	})

	t.Run("constructor with error", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register
		assert.Nil(manioc.RegisterConstructor[IMyService](NewMyServiceWithError, manioc.WithContainer(ctr)))

		// resolve
		ret, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Nil(err)

		// verify that the value of ret is the instance of MyService
		_, ok := ret.(*MyService)
		assert.True(ok)
	})

	t.Run("error propagation from constructor", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()

		// register
		assert.Nil(manioc.RegisterConstructor[IMyService](NewMyServiceRaiseError, manioc.WithContainer(ctr)))

		// resolve
		ret, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
		assert.Nil(ret)
		assert.Error(err)
	})
}

func Test_RegisterInstanceAndResolve(t *testing.T) {
	assert := assert.New(t)

	ctr := manioc.NewContainer()

	// nil instance is not allowed
	assert.Error(manioc.RegisterInstance[IMyService](nil, manioc.WithContainer(ctr)))

	// register
	instance := &MyService{}
	assert.Nil(manioc.RegisterInstance[IMyService](instance, manioc.WithContainer(ctr)))

	// resolve
	ret, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)

	// verify that the value of ret is the instance of MyService
	value, ok := ret.(*MyService)
	assert.True(ok)

	// when using RegisterInstance, The resolved instance matches the registered one.
	assert.Same(instance, value)
}

func Test_RegisterAndResolveWithKey(t *testing.T) {
	t.Run("register with anonymous key", func(t *testing.T) {
		assert := assert.New(t)

		// register
		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IMyService, MyService](manioc.WithContainer(ctr)))

		// resolve
		ret, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Nil(err)

		// since MyService is registered with anonymous key,
		// it cannot be resolved with any other keys
		_, err = manioc.Resolve[IMyService](
			manioc.WithScope(ctr),
			manioc.WithResolveKey("mykey"),
		)
		assert.Error(err)
	})

	t.Run("register with key", func(t *testing.T) {
		assert := assert.New(t)

		// register with key
		ctr := manioc.NewContainer()
		key := "mykey"
		assert.Nil(manioc.Register[IMyService, MyService](
			manioc.WithContainer(ctr),
			manioc.WithRegisterKey(key),
		))

		// resolve
		ret, err := manioc.Resolve[IMyService](
			manioc.WithScope(ctr),
			manioc.WithResolveKey(key),
		)
		assert.NotNil(ret)
		assert.Nil(err)

		// since MyService is registered with key,
		// it cannot be resolved with anonymous key
		_, err = manioc.Resolve[IMyService](manioc.WithScope(ctr))
		assert.Error(err)
	})
}

func Test_Register_Errors(t *testing.T) {
	t.Run("resolve not registered service", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		_, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
		assert.Error(err)
	})

	t.Run("runtime check: TImplementation should implement TInterface", func(t *testing.T) {
		// We want to constrain the TImplementation type by TInterface
		// (i.e. something like Register[TInterface any, TImplementation TInterface](...)),
		// but it is not supported in Golang 1.18.
		// We don't know if this will ever be statically checkable,
		// but until then, anyway, we decided to implement it with reflection.
		// It is designed to cause panic if verification fails.
		ctr := manioc.NewContainer()
		assert.Panics(t, func() { _ = manioc.Register[IMyService, bool](manioc.WithContainer(ctr)) })
		assert.Panics(t, func() { _ = manioc.Register[IMyService, int](manioc.WithContainer(ctr)) })
		assert.Panics(t, func() { _ = manioc.Register[IMyService, string](manioc.WithContainer(ctr)) })
		assert.Panics(t, func() { _ = manioc.Register[IMyService, *struct{}](manioc.WithContainer(ctr)) })
		assert.NotPanics(t, func() { _ = manioc.Register[IMyService, MyService](manioc.WithContainer(ctr)) })
	})
}

func Test_RegisterConstructor_Errors(t *testing.T) {
	t.Run("runtime check: TConstructor should return TInterface", func(t *testing.T) {
		// We want to constrain the TConstructor type to be any function with a return value of type TInterface.
		// As of Golang 1.18, this seems to be difficult to achieve.
		// Suggestions:
		// 1. The type `func (...any) T` accepts any function whose return type is `T`.
		//    - It is not allowed in Golang, but, for example, is allowed in TypeScript.
		//      ```ts
		//      type AnyFuncReturns<T> = (...args: any[]) => T
		//      const fn1: AnyFuncReturns<number> = function (): number { return 42 }
		//      const fn2: AnyFuncReturns<number> = function (value: string): number { return 42 }
		//      const fn3: AnyFuncReturns<number> = function (value: string, options: any): number { return 42 }
		//      ```
		//    - We feel that this extension is reasonable because it does not require a strong extension of generics feature
		//      and keeps backward compatibility within the language specification.
		// 2. Introducing a generic type `Func[T]` that represents all of the functions whose return type is `T`.
		//    - The advantage of this is that only the generics feature needs to be extended,
		//      but a built-in special type must be provided.
		// 3. Introducing "variadic type parameters".
		//    - Something like `RegisterConstructor[TInterface any, ...TArgs](ctor Func[TArgs..., TInterface])`.
		//    - This is a more powerful generics extension than 2. It would be interesting to have,
		//      but if I were a golang developer, I would choose another plan.
		// Anyway, we decided to implement it with reflection.
		// It is designed to cause panic if verification fails.
		ctr := manioc.NewContainer()
		assert.Panics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](42, manioc.WithContainer(ctr))
		})
		assert.Panics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](struct{}{}, manioc.WithContainer(ctr))
		})
		assert.Panics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](&struct{}{}, manioc.WithContainer(ctr))
		})
		assert.Panics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](func() {}, manioc.WithContainer(ctr))
		})
		assert.Panics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](func() bool { return false }, manioc.WithContainer(ctr))
		})
		assert.Panics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](func() int { return 42 }, manioc.WithContainer(ctr))
		})
		assert.Panics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](func() string { return "hello world" }, manioc.WithContainer(ctr))
		})
		assert.Panics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](func() *struct{} { return &struct{}{} }, manioc.WithContainer(ctr))
		})
		assert.Panics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](func() (IMyService, int) { return nil, 42 }, manioc.WithContainer(ctr))
		})
		assert.Panics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](func() (int, error) { return 42, nil }, manioc.WithContainer(ctr))
		})
		assert.NotPanics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](NewMyService, manioc.WithContainer(ctr))
		})
		assert.NotPanics(t, func() {
			_ = manioc.RegisterConstructor[IMyService](
				func(a int, b string, c any) IMyService { return nil },
				manioc.WithContainer(ctr),
			)
		})
	})

	t.Run("ctor should not be nil", func(t *testing.T) {
		ctr := manioc.NewContainer()
		assert.Error(t, manioc.RegisterConstructor[IMyService, func() IMyService](nil, manioc.WithContainer(ctr)))
	})
}

func Test_MustResolve(t *testing.T) {
	t.Run("MustResolve causes panic for resolution errors", func(t *testing.T) {
		assert.Panics(t, func() {
			ctr := manioc.NewContainer()
			_ = manioc.MustResolve[IMyService](manioc.WithScope(ctr))
		})
	})

	t.Run("MustResolve returns the instance when resolution is successfully finished", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[IMyService, MyService](manioc.WithContainer(ctr)))

		ret := manioc.MustResolve[IMyService](manioc.WithScope(ctr))
		assert.NotNil(ret)
		_, ok := ret.(*MyService)
		assert.True(ok)
	})
}
