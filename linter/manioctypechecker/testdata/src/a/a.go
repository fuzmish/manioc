package a

import (
	"github.com/fuzmish/manioc"
)

// valid
var v1 = manioc.Register[IMyService, MyService]
var v2 = manioc.Register[IMyService, *MyService]
var v3 = manioc.RegisterConstructor[IMyService, func() *MyService]
var v4 = manioc.RegisterConstructor[IMyService, func() (*MyService, error)]
var v5 = manioc.Register[int, int]
var v6 = manioc.Register[*int, *int]
var v7 = manioc.Register[MyService, MyService]
var v8 = manioc.Register[*MyService, *MyService]

// invalid
var i1 = manioc.Register[IMyService, struct{}]                                   // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`a\\.IMyService`"
var i2 = manioc.Register[IMyService, *struct{}]                                  // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`a\\.IMyService`"
var i3 = manioc.RegisterConstructor[IMyService, func()]                          // want "The number of TConstructor return values should be either one or two"
var i4 = manioc.RegisterConstructor[IMyService, func() MyService]                // want "The first return type of TConstructor `a\\.MyService` is not assignable to `a.IMyService`"
var i5 = manioc.RegisterConstructor[IMyService, func() (*MyService, int)]        // want "The second return type of TConstructor should be `error`, but `int` is given"
var i6 = manioc.RegisterConstructor[IMyService, func() (MyService, error)]       // want "The first return type of TConstructor `a\\.MyService` is not assignable to `a\\.IMyService`"
var i7 = manioc.RegisterConstructor[IMyService, func() (*MyService, error, int)] // want "The number of TConstructor return values should be either one or two"
var i8 = manioc.Register[int, string]                                            // want "TImplementation=`string` is not assignable to TInterface=`int`"
var i9 = manioc.Register[*int, *string]                                          // want "TImplementation=`\\*string` is not assignable to TInterface=`\\*int`"
var i10 = manioc.Register[MyService, struct{}]                                   // want "TImplementation=`struct{}` is not assignable to TInterface=`a\\.MyService`"
var i11 = manioc.Register[*MyService, *struct{}]                                 // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`\\*a\\.MyService`"

type IMyService interface {
	doSomething()
}

type MyService struct {
	value int
}

func (s *MyService) doSomething() {}

func NewMyService() *MyService {
	return &MyService{}
}

func NewMyServiceWithError() (*MyService, error) {
	return &MyService{}, nil
}

func valid() {
	ctr := manioc.NewContainer()

	// v1
	_ = manioc.Register[IMyService, MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterSingleton[IMyService, MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterScoped[IMyService, MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterTransient[IMyService, MyService](manioc.WithContainer(ctr))

	// v2
	_ = manioc.Register[IMyService, *MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterSingleton[IMyService, *MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterScoped[IMyService, *MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterTransient[IMyService, *MyService](manioc.WithContainer(ctr))

	// v3
	_ = manioc.RegisterConstructor[IMyService](NewMyService, manioc.WithContainer(ctr))
	_ = manioc.RegisterSingletonConstructor[IMyService](NewMyService, manioc.WithContainer(ctr))
	_ = manioc.RegisterScopedConstructor[IMyService](NewMyService, manioc.WithContainer(ctr))
	_ = manioc.RegisterTransientConstructor[IMyService](NewMyService, manioc.WithContainer(ctr))

	// v4
	_ = manioc.RegisterConstructor[IMyService](NewMyServiceWithError, manioc.WithContainer(ctr))
	_ = manioc.RegisterSingletonConstructor[IMyService](NewMyServiceWithError, manioc.WithContainer(ctr))
	_ = manioc.RegisterScopedConstructor[IMyService](NewMyServiceWithError, manioc.WithContainer(ctr))
	_ = manioc.RegisterTransientConstructor[IMyService](NewMyServiceWithError, manioc.WithContainer(ctr))

	// v5
	_ = manioc.Register[int, int](manioc.WithContainer(ctr))
	_ = manioc.RegisterSingleton[int, int](manioc.WithContainer(ctr))
	_ = manioc.RegisterScoped[int, int](manioc.WithContainer(ctr))
	_ = manioc.Register[int, int](manioc.WithContainer(ctr))

	// v6
	_ = manioc.Register[*int, *int](manioc.WithContainer(ctr))
	_ = manioc.RegisterSingleton[*int, *int](manioc.WithContainer(ctr))
	_ = manioc.RegisterScoped[*int, *int](manioc.WithContainer(ctr))
	_ = manioc.RegisterTransient[*int, *int](manioc.WithContainer(ctr))

	// v7
	_ = manioc.Register[MyService, MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterSingleton[MyService, MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterScoped[MyService, MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterTransient[MyService, MyService](manioc.WithContainer(ctr))

	// v8
	_ = manioc.Register[*MyService, *MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterSingleton[*MyService, *MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterScoped[*MyService, *MyService](manioc.WithContainer(ctr))
	_ = manioc.RegisterTransient[*MyService, *MyService](manioc.WithContainer(ctr))
}

func invalid() {
	ctr := manioc.NewContainer()

	// i1
	_ = manioc.Register[IMyService, struct{}](manioc.WithContainer(ctr))          // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`a\\.IMyService`"
	_ = manioc.RegisterSingleton[IMyService, struct{}](manioc.WithContainer(ctr)) // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`a\\.IMyService`"
	_ = manioc.RegisterScoped[IMyService, struct{}](manioc.WithContainer(ctr))    // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`a\\.IMyService`"
	_ = manioc.RegisterTransient[IMyService, struct{}](manioc.WithContainer(ctr)) // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`a\\.IMyService`"

	// i2
	_ = manioc.Register[IMyService, *struct{}](manioc.WithContainer(ctr))          // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`a\\.IMyService`"
	_ = manioc.RegisterSingleton[IMyService, *struct{}](manioc.WithContainer(ctr)) // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`a\\.IMyService`"
	_ = manioc.RegisterScoped[IMyService, *struct{}](manioc.WithContainer(ctr))    // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`a\\.IMyService`"
	_ = manioc.RegisterTransient[IMyService, *struct{}](manioc.WithContainer(ctr)) // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`a\\.IMyService`"

	// i3
	_ = manioc.RegisterConstructor[IMyService](func() {}, manioc.WithContainer(ctr))          // want "The number of TConstructor return values should be either one or two"
	_ = manioc.RegisterSingletonConstructor[IMyService](func() {}, manioc.WithContainer(ctr)) // want "The number of TConstructor return values should be either one or two"
	_ = manioc.RegisterScopedConstructor[IMyService](func() {}, manioc.WithContainer(ctr))    // want "The number of TConstructor return values should be either one or two"
	_ = manioc.RegisterTransientConstructor[IMyService](func() {}, manioc.WithContainer(ctr)) // want "The number of TConstructor return values should be either one or two"

	// i4
	_ = manioc.RegisterConstructor[IMyService](func() MyService { return MyService{} }, manioc.WithContainer(ctr))          // want "The first return type of TConstructor `a\\.MyService` is not assignable to `a.IMyService`"
	_ = manioc.RegisterSingletonConstructor[IMyService](func() MyService { return MyService{} }, manioc.WithContainer(ctr)) // want "The first return type of TConstructor `a\\.MyService` is not assignable to `a.IMyService`"
	_ = manioc.RegisterScopedConstructor[IMyService](func() MyService { return MyService{} }, manioc.WithContainer(ctr))    // want "The first return type of TConstructor `a\\.MyService` is not assignable to `a.IMyService`"
	_ = manioc.RegisterTransientConstructor[IMyService](func() MyService { return MyService{} }, manioc.WithContainer(ctr)) // want "The first return type of TConstructor `a\\.MyService` is not assignable to `a.IMyService`"

	// i5
	_ = manioc.RegisterConstructor[IMyService](func() (*MyService, int) { return &MyService{}, 42 }, manioc.WithContainer(ctr))          // want "The second return type of TConstructor should be `error`, but `int` is given"
	_ = manioc.RegisterSingletonConstructor[IMyService](func() (*MyService, int) { return &MyService{}, 42 }, manioc.WithContainer(ctr)) // want "The second return type of TConstructor should be `error`, but `int` is given"
	_ = manioc.RegisterScopedConstructor[IMyService](func() (*MyService, int) { return &MyService{}, 42 }, manioc.WithContainer(ctr))    // want "The second return type of TConstructor should be `error`, but `int` is given"
	_ = manioc.RegisterTransientConstructor[IMyService](func() (*MyService, int) { return &MyService{}, 42 }, manioc.WithContainer(ctr)) // want "The second return type of TConstructor should be `error`, but `int` is given"

	// i6
	_ = manioc.RegisterConstructor[IMyService](func() (MyService, error) { return MyService{}, nil }, manioc.WithContainer(ctr))          // want "The first return type of TConstructor `a\\.MyService` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterSingletonConstructor[IMyService](func() (MyService, error) { return MyService{}, nil }, manioc.WithContainer(ctr)) // want "The first return type of TConstructor `a\\.MyService` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterScopedConstructor[IMyService](func() (MyService, error) { return MyService{}, nil }, manioc.WithContainer(ctr))    // want "The first return type of TConstructor `a\\.MyService` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterTransientConstructor[IMyService](func() (MyService, error) { return MyService{}, nil }, manioc.WithContainer(ctr)) // want "The first return type of TConstructor `a\\.MyService` is not assignable to `a\\.IMyService`"

	// i7
	_ = manioc.RegisterConstructor[IMyService](func() (*MyService, error, int) { return &MyService{}, nil, 42 }, manioc.WithContainer(ctr))          // want "The number of TConstructor return values should be either one or two"
	_ = manioc.RegisterSingletonConstructor[IMyService](func() (*MyService, error, int) { return &MyService{}, nil, 42 }, manioc.WithContainer(ctr)) // want "The number of TConstructor return values should be either one or two"
	_ = manioc.RegisterScopedConstructor[IMyService](func() (*MyService, error, int) { return &MyService{}, nil, 42 }, manioc.WithContainer(ctr))    // want "The number of TConstructor return values should be either one or two"
	_ = manioc.RegisterTransientConstructor[IMyService](func() (*MyService, error, int) { return &MyService{}, nil, 42 }, manioc.WithContainer(ctr)) // want "The number of TConstructor return values should be either one or two"

	// i8
	_ = manioc.Register[int, string](manioc.WithContainer(ctr))          // want "TImplementation=`string` is not assignable to TInterface=`int`"
	_ = manioc.RegisterSingleton[int, string](manioc.WithContainer(ctr)) // want "TImplementation=`string` is not assignable to TInterface=`int`"
	_ = manioc.RegisterScoped[int, string](manioc.WithContainer(ctr))    // want "TImplementation=`string` is not assignable to TInterface=`int`"
	_ = manioc.RegisterTransient[int, string](manioc.WithContainer(ctr)) // want "TImplementation=`string` is not assignable to TInterface=`int`"

	// i9
	_ = manioc.Register[*int, *string](manioc.WithContainer(ctr))          // want "TImplementation=`\\*string` is not assignable to TInterface=`\\*int`"
	_ = manioc.RegisterSingleton[*int, *string](manioc.WithContainer(ctr)) // want "TImplementation=`\\*string` is not assignable to TInterface=`\\*int`"
	_ = manioc.RegisterScoped[*int, *string](manioc.WithContainer(ctr))    // want "TImplementation=`\\*string` is not assignable to TInterface=`\\*int`"
	_ = manioc.RegisterTransient[*int, *string](manioc.WithContainer(ctr)) // want "TImplementation=`\\*string` is not assignable to TInterface=`\\*int`"

	// i10
	_ = manioc.Register[MyService, struct{}](manioc.WithContainer(ctr))          // want "TImplementation=`struct{}` is not assignable to TInterface=`a\\.MyService`"
	_ = manioc.RegisterSingleton[MyService, struct{}](manioc.WithContainer(ctr)) // want "TImplementation=`struct{}` is not assignable to TInterface=`a\\.MyService`"
	_ = manioc.RegisterScoped[MyService, struct{}](manioc.WithContainer(ctr))    // want "TImplementation=`struct{}` is not assignable to TInterface=`a\\.MyService`"
	_ = manioc.RegisterTransient[MyService, struct{}](manioc.WithContainer(ctr)) // want "TImplementation=`struct{}` is not assignable to TInterface=`a\\.MyService`"

	// i11
	_ = manioc.Register[*MyService, *struct{}](manioc.WithContainer(ctr))          // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`\\*a\\.MyService`"
	_ = manioc.RegisterSingleton[*MyService, *struct{}](manioc.WithContainer(ctr)) // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`\\*a\\.MyService`"
	_ = manioc.RegisterScoped[*MyService, *struct{}](manioc.WithContainer(ctr))    // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`\\*a\\.MyService`"
	_ = manioc.RegisterTransient[*MyService, *struct{}](manioc.WithContainer(ctr)) // want "TImplementation=`\\*struct{}` is not assignable to TInterface=`\\*a\\.MyService`"
}
