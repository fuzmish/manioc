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
var v9 = manioc.ResolveFunction[IMyService, func() *MyService]
var v10 = manioc.ResolveFunction[IMyService, func() (*MyService, error)]
var v11 = manioc.MustResolveFunction[IMyService, func() *MyService]
var v12 = manioc.MustResolveFunction[IMyService, func() (*MyService, error)]

// invalid
var i1 = manioc.Register[IMyService, struct{}]                                    // want "`struct\\{\\}` is not assignable to `a\\.IMyService`"
var i2 = manioc.Register[IMyService, *struct{}]                                   // want "`\\*struct\\{\\}` is not assignable to `a\\.IMyService`"
var i3 = manioc.Register[IMyService, IMyService]                                  // want "The implementation type `a\\.IMyService` should not be an interface"
var i4 = manioc.Register[any, any]                                                // want "The implementation type `any` should not be an interface"
var i5 = manioc.RegisterConstructor[IMyService, func()]                           // want "The number of function return values should be either one or two"
var i6 = manioc.RegisterConstructor[IMyService, func() MyService]                 // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
var i7 = manioc.RegisterConstructor[IMyService, func() (*MyService, int)]         // want "The type of the second return value should be `error`, but `int` is given"
var i8 = manioc.RegisterConstructor[IMyService, func() (MyService, error)]        // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
var i9 = manioc.RegisterConstructor[IMyService, func() (*MyService, error, int)]  // want "The number of function return values should be either one or two"
var i10 = manioc.Register[int, string]                                            // want "`string` is not assignable to `int`"
var i11 = manioc.Register[*int, *string]                                          // want "`\\*string` is not assignable to `\\*int`"
var i12 = manioc.Register[MyService, struct{}]                                    // want "`struct\\{\\}` is not assignable to `a\\.MyService`"
var i13 = manioc.Register[*MyService, *struct{}]                                  // want "`\\*struct\\{\\}` is not assignable to `\\*a\\.MyService`"
var i14 = manioc.ResolveFunction[IMyService, func()]                              // want "The number of function return values should be either one or two"
var i15 = manioc.ResolveFunction[IMyService, func() MyService]                    // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
var i16 = manioc.ResolveFunction[IMyService, func() (*MyService, int)]            // want "The type of the second return value should be `error`, but `int` is given"
var i17 = manioc.ResolveFunction[IMyService, func() (MyService, error)]           // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
var i18 = manioc.ResolveFunction[IMyService, func() (*MyService, error, int)]     // want "The number of function return values should be either one or two"
var i19 = manioc.MustResolveFunction[IMyService, func()]                          // want "The number of function return values should be either one or two"
var i20 = manioc.MustResolveFunction[IMyService, func() MyService]                // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
var i21 = manioc.MustResolveFunction[IMyService, func() (*MyService, int)]        // want "The type of the second return value should be `error`, but `int` is given"
var i22 = manioc.MustResolveFunction[IMyService, func() (MyService, error)]       // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
var i23 = manioc.MustResolveFunction[IMyService, func() (*MyService, error, int)] // want "The number of function return values should be either one or two"

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
	// v1
	_ = manioc.Register[IMyService, MyService]()
	_ = manioc.RegisterSingleton[IMyService, MyService]()
	_ = manioc.RegisterScoped[IMyService, MyService]()
	_ = manioc.RegisterTransient[IMyService, MyService]()

	// v2
	_ = manioc.Register[IMyService, *MyService]()
	_ = manioc.RegisterSingleton[IMyService, *MyService]()
	_ = manioc.RegisterScoped[IMyService, *MyService]()
	_ = manioc.RegisterTransient[IMyService, *MyService]()

	// v3
	_ = manioc.RegisterConstructor[IMyService](NewMyService)
	_ = manioc.RegisterSingletonConstructor[IMyService](NewMyService)
	_ = manioc.RegisterScopedConstructor[IMyService](NewMyService)
	_ = manioc.RegisterTransientConstructor[IMyService](NewMyService)

	// v4
	_ = manioc.RegisterConstructor[IMyService](NewMyServiceWithError)
	_ = manioc.RegisterSingletonConstructor[IMyService](NewMyServiceWithError)
	_ = manioc.RegisterScopedConstructor[IMyService](NewMyServiceWithError)
	_ = manioc.RegisterTransientConstructor[IMyService](NewMyServiceWithError)

	// v5
	_ = manioc.Register[int, int]()
	_ = manioc.RegisterSingleton[int, int]()
	_ = manioc.RegisterScoped[int, int]()
	_ = manioc.Register[int, int]()

	// v6
	_ = manioc.Register[*int, *int]()
	_ = manioc.RegisterSingleton[*int, *int]()
	_ = manioc.RegisterScoped[*int, *int]()
	_ = manioc.RegisterTransient[*int, *int]()

	// v7
	_ = manioc.Register[MyService, MyService]()
	_ = manioc.RegisterSingleton[MyService, MyService]()
	_ = manioc.RegisterScoped[MyService, MyService]()
	_ = manioc.RegisterTransient[MyService, MyService]()

	// v8
	_ = manioc.Register[*MyService, *MyService]()
	_ = manioc.RegisterSingleton[*MyService, *MyService]()
	_ = manioc.RegisterScoped[*MyService, *MyService]()
	_ = manioc.RegisterTransient[*MyService, *MyService]()

	// v9
	_, _ = manioc.ResolveFunction[IMyService](func() *MyService { return &MyService{} })

	// v10
	_, _ = manioc.ResolveFunction[IMyService](func() (*MyService, error) { return &MyService{}, nil })

	// v11
	_ = manioc.MustResolveFunction[IMyService](func() *MyService { return &MyService{} })

	// v12
	_ = manioc.MustResolveFunction[IMyService](func() (*MyService, error) { return &MyService{}, nil })
}

func invalid() {
	// i1
	_ = manioc.Register[IMyService, struct{}]()          // want "`struct\\{\\}` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterSingleton[IMyService, struct{}]() // want "`struct\\{\\}` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterScoped[IMyService, struct{}]()    // want "`struct\\{\\}` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterTransient[IMyService, struct{}]() // want "`struct\\{\\}` is not assignable to `a\\.IMyService`"

	// i2
	_ = manioc.Register[IMyService, *struct{}]()          // want "`\\*struct\\{\\}` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterSingleton[IMyService, *struct{}]() // want "`\\*struct\\{\\}` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterScoped[IMyService, *struct{}]()    // want "`\\*struct\\{\\}` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterTransient[IMyService, *struct{}]() // want "`\\*struct\\{\\}` is not assignable to `a\\.IMyService`"

	// i3
	_ = manioc.Register[IMyService, IMyService]()          // want "The implementation type `a\\.IMyService` should not be an interface"
	_ = manioc.RegisterSingleton[IMyService, IMyService]() // want "The implementation type `a\\.IMyService` should not be an interface"
	_ = manioc.RegisterScoped[IMyService, IMyService]()    // want "The implementation type `a\\.IMyService` should not be an interface"
	_ = manioc.RegisterTransient[IMyService, IMyService]() // want "The implementation type `a\\.IMyService` should not be an interface"

	// i4
	_ = manioc.Register[any, any]()          // want "The implementation type `any` should not be an interface"
	_ = manioc.RegisterSingleton[any, any]() // want "The implementation type `any` should not be an interface"
	_ = manioc.RegisterScoped[any, any]()    // want "The implementation type `any` should not be an interface"
	_ = manioc.RegisterTransient[any, any]() // want "The implementation type `any` should not be an interface"

	// i5
	_ = manioc.RegisterConstructor[IMyService](func() {})          // want "The number of function return values should be either one or two"
	_ = manioc.RegisterSingletonConstructor[IMyService](func() {}) // want "The number of function return values should be either one or two"
	_ = manioc.RegisterScopedConstructor[IMyService](func() {})    // want "The number of function return values should be either one or two"
	_ = manioc.RegisterTransientConstructor[IMyService](func() {}) // want "The number of function return values should be either one or two"

	// i6
	_ = manioc.RegisterConstructor[IMyService](func() MyService { return MyService{} })          // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterSingletonConstructor[IMyService](func() MyService { return MyService{} }) // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterScopedConstructor[IMyService](func() MyService { return MyService{} })    // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterTransientConstructor[IMyService](func() MyService { return MyService{} }) // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"

	// i7
	_ = manioc.RegisterConstructor[IMyService](func() (*MyService, int) { return &MyService{}, 42 })          // want "The type of the second return value should be `error`, but `int` is given"
	_ = manioc.RegisterSingletonConstructor[IMyService](func() (*MyService, int) { return &MyService{}, 42 }) // want "The type of the second return value should be `error`, but `int` is given"
	_ = manioc.RegisterScopedConstructor[IMyService](func() (*MyService, int) { return &MyService{}, 42 })    // want "The type of the second return value should be `error`, but `int` is given"
	_ = manioc.RegisterTransientConstructor[IMyService](func() (*MyService, int) { return &MyService{}, 42 }) // want "The type of the second return value should be `error`, but `int` is given"

	// i8
	_ = manioc.RegisterConstructor[IMyService](func() (MyService, error) { return MyService{}, nil })          // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterSingletonConstructor[IMyService](func() (MyService, error) { return MyService{}, nil }) // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterScopedConstructor[IMyService](func() (MyService, error) { return MyService{}, nil })    // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"
	_ = manioc.RegisterTransientConstructor[IMyService](func() (MyService, error) { return MyService{}, nil }) // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"

	// i9
	_ = manioc.RegisterConstructor[IMyService](func() (*MyService, error, int) { return &MyService{}, nil, 42 })          // want "The number of function return values should be either one or two"
	_ = manioc.RegisterSingletonConstructor[IMyService](func() (*MyService, error, int) { return &MyService{}, nil, 42 }) // want "The number of function return values should be either one or two"
	_ = manioc.RegisterScopedConstructor[IMyService](func() (*MyService, error, int) { return &MyService{}, nil, 42 })    // want "The number of function return values should be either one or two"
	_ = manioc.RegisterTransientConstructor[IMyService](func() (*MyService, error, int) { return &MyService{}, nil, 42 }) // want "The number of function return values should be either one or two"

	// i10
	_ = manioc.Register[int, string]()          // want "`string` is not assignable to `int`"
	_ = manioc.RegisterSingleton[int, string]() // want "`string` is not assignable to `int`"
	_ = manioc.RegisterScoped[int, string]()    // want "`string` is not assignable to `int`"
	_ = manioc.RegisterTransient[int, string]() // want "`string` is not assignable to `int`"

	// i11
	_ = manioc.Register[*int, *string]()          // want "`\\*string` is not assignable to `\\*int`"
	_ = manioc.RegisterSingleton[*int, *string]() // want "`\\*string` is not assignable to `\\*int`"
	_ = manioc.RegisterScoped[*int, *string]()    // want "`\\*string` is not assignable to `\\*int`"
	_ = manioc.RegisterTransient[*int, *string]() // want "`\\*string` is not assignable to `\\*int`"

	// i12
	_ = manioc.Register[MyService, struct{}]()          // want "`struct\\{\\}` is not assignable to `a\\.MyService`"
	_ = manioc.RegisterSingleton[MyService, struct{}]() // want "`struct\\{\\}` is not assignable to `a\\.MyService`"
	_ = manioc.RegisterScoped[MyService, struct{}]()    // want "`struct\\{\\}` is not assignable to `a\\.MyService`"
	_ = manioc.RegisterTransient[MyService, struct{}]() // want "`struct\\{\\}` is not assignable to `a\\.MyService`"

	// i13
	_ = manioc.Register[*MyService, *struct{}]()          // want "`\\*struct\\{\\}` is not assignable to `\\*a\\.MyService`"
	_ = manioc.RegisterSingleton[*MyService, *struct{}]() // want "`\\*struct\\{\\}` is not assignable to `\\*a\\.MyService`"
	_ = manioc.RegisterScoped[*MyService, *struct{}]()    // want "`\\*struct\\{\\}` is not assignable to `\\*a\\.MyService`"
	_ = manioc.RegisterTransient[*MyService, *struct{}]() // want "`\\*struct\\{\\}` is not assignable to `\\*a\\.MyService`"

	// i14
	_, _ = manioc.ResolveFunction[IMyService](func() {}) // want "The number of function return values should be either one or two"

	// i15
	_, _ = manioc.ResolveFunction[IMyService](func() MyService { return MyService{} }) // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"

	// i16
	_, _ = manioc.ResolveFunction[IMyService](func() (*MyService, int) { return &MyService{}, 42 }) // want "The type of the second return value should be `error`, but `int` is given"

	// i17
	_, _ = manioc.ResolveFunction[IMyService](func() (MyService, error) { return MyService{}, nil }) // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"

	// i18
	_, _ = manioc.ResolveFunction[IMyService](func() (*MyService, error, int) { return &MyService{}, nil, 42 }) // want "The number of function return values should be either one or two"

	// i19
	_ = manioc.MustResolveFunction[IMyService](func() {}) // want "The number of function return values should be either one or two"

	// i20
	_ = manioc.MustResolveFunction[IMyService](func() MyService { return MyService{} }) // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"

	// i21
	_ = manioc.MustResolveFunction[IMyService](func() (*MyService, int) { return &MyService{}, 42 }) // want "The type of the second return value should be `error`, but `int` is given"

	// i22
	_ = manioc.MustResolveFunction[IMyService](func() (MyService, error) { return MyService{}, nil }) // want "The type of the first return value `a\\.MyService` is not assignable to `a\\.IMyService`"

	// i23
	_ = manioc.MustResolveFunction[IMyService](func() (*MyService, error, int) { return &MyService{}, nil, 42 }) // want "The number of function return values should be either one or two"
}
