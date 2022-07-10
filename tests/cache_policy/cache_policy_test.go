package manioc_cache_policy_test

import (
	"testing"

	"github.com/fuzmish/manioc"
	"github.com/stretchr/testify/assert"
)

type IMyService interface {
	doSomething()
}

// MyService implements IMyService
type MyService struct {
	Value int
}

func (s *MyService) doSomething() {}

func NewMyService() *MyService {
	return &MyService{Value: 42}
}

func Test_CachePolicy_NeverCache(t *testing.T) {
	registerFunctions := map[string]func(ctr manioc.Container) error{
		"Register: the default cache policy is NeverCache": func(ctr manioc.Container) error {
			return manioc.Register[IMyService, MyService](manioc.WithContainer(ctr))
		},
		"Register: set NeverCache policy explicitly": func(ctr manioc.Container) error {
			return manioc.Register[IMyService, MyService](
				manioc.WithContainer(ctr),
				manioc.WithCachePolicy(manioc.NeverCache),
			)
		},
		"RegisterTransient: the helper function": func(ctr manioc.Container) error {
			return manioc.RegisterTransient[IMyService, MyService](manioc.WithContainer(ctr))
		},
		"RegisterConstructor: the default cache policy is NeverCache": func(ctr manioc.Container) error {
			return manioc.RegisterConstructor[IMyService](NewMyService, manioc.WithContainer(ctr))
		},
		"RegisterConstructor: set NeverCache policy explicitly": func(ctr manioc.Container) error {
			return manioc.RegisterConstructor[IMyService](
				NewMyService,
				manioc.WithContainer(ctr),
				manioc.WithCachePolicy(manioc.NeverCache),
			)
		},
		"RegisterTransientConstructor: the helper function": func(ctr manioc.Container) error {
			return manioc.RegisterTransientConstructor[IMyService](NewMyService, manioc.WithContainer(ctr))
		},
	}
	for name, registerFunction := range registerFunctions {
		registerFunction := registerFunction
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			ctr := manioc.NewContainer()
			assert.Nil(registerFunction(ctr))

			ret1, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
			assert.Nil(err)
			assert.NotNil(ret1)
			ret2, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
			assert.Nil(err)
			assert.NotNil(ret2)

			// Under NeverCache policy, a new instance is created for each resolution.
			assert.NotSame(ret1, ret2)
		})
	}
}

func Test_CachePolicy_ScopedCache(t *testing.T) {
	registerFunctions := map[string]func(ctr manioc.Container) error{
		"Register: set ScopedCache policy explicitly": func(ctr manioc.Container) error {
			return manioc.Register[IMyService, MyService](
				manioc.WithContainer(ctr),
				manioc.WithCachePolicy(manioc.ScopedCache),
			)
		},
		"RegisterScoped: the helper function": func(ctr manioc.Container) error {
			return manioc.RegisterScoped[IMyService, MyService](manioc.WithContainer(ctr))
		},
		"RegisterConstructor: set ScopedCache policy explicitly": func(ctr manioc.Container) error {
			return manioc.RegisterConstructor[IMyService](
				NewMyService,
				manioc.WithContainer(ctr),
				manioc.WithCachePolicy(manioc.ScopedCache),
			)
		},
		"RegisterScopedConstructor: the helper function": func(ctr manioc.Container) error {
			return manioc.RegisterScopedConstructor[IMyService](NewMyService, manioc.WithContainer(ctr))
		},
	}
	for name, registerFunction := range registerFunctions {
		registerFunction := registerFunction
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			ctr := manioc.NewContainer()
			assert.Nil(registerFunction(ctr))

			// Open scope.
			scope1, _ := ctr.OpenScope()

			ret1, err := manioc.Resolve[IMyService](manioc.WithScope(scope1))
			assert.Nil(err)
			assert.NotNil(ret1)
			ret2, err := manioc.Resolve[IMyService](manioc.WithScope(scope1))
			assert.Nil(err)
			assert.NotNil(ret2)

			// Under ScopedCache policy, an instance is created only once and will be reused within the same scope.
			assert.Same(ret1, ret2)

			// Open another scope.
			scope2, _ := ctr.OpenScope()

			ret3, err := manioc.Resolve[IMyService](manioc.WithScope(scope2))
			assert.Nil(err)
			assert.NotNil(ret3)

			// For different scopes, an instance is created separately.
			assert.NotSame(ret1, ret3)
		})
	}
}

func Test_CachePolicy_GlobalCache(t *testing.T) {
	registerFunctions := map[string]func(ctr manioc.Container) error{
		"Register: set GlobalCache policy explicitly": func(ctr manioc.Container) error {
			return manioc.Register[IMyService, MyService](
				manioc.WithContainer(ctr),
				manioc.WithCachePolicy(manioc.GlobalCache),
			)
		},
		"RegisterSingleton: the helper function": func(ctr manioc.Container) error {
			return manioc.RegisterSingleton[IMyService, MyService](manioc.WithContainer(ctr))
		},
		"RegisterConstructor: set GlobalCache policy explicitly": func(ctr manioc.Container) error {
			return manioc.RegisterConstructor[IMyService](
				NewMyService,
				manioc.WithContainer(ctr),
				manioc.WithCachePolicy(manioc.GlobalCache),
			)
		},
		"RegisterTransientConstructor: the helper function": func(ctr manioc.Container) error {
			return manioc.RegisterSingletonConstructor[IMyService](NewMyService, manioc.WithContainer(ctr))
		},
		"RegisterInstance: the default cache policy is NeverCache": func(ctr manioc.Container) error {
			return manioc.RegisterInstance[IMyService](&MyService{Value: 42}, manioc.WithContainer(ctr))
		},
	}
	for name, registerFunction := range registerFunctions {
		registerFunction := registerFunction
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			ctr := manioc.NewContainer()
			assert.Nil(registerFunction(ctr))

			ret1, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
			assert.Nil(err)
			assert.NotNil(ret1)
			ret2, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
			assert.Nil(err)
			assert.NotNil(ret2)

			// Under GlobalCache policy, an instance is created only once and will be reused.
			assert.Same(ret1, ret2)
		})
	}
}

func Test_CachePolicy_On_RegisterInstance(t *testing.T) {
	// For RegisterInstance, only GlobalCache policy is available.
	// If you specify the other policy, it will be ignored.
	// It would be nice to set this function so that it does not accept the CachePolicy option,
	// but it is tedious to realize in FOP.
	// This can be solved when golang introduces optional arguments.
	policies := map[string]manioc.CachePolicy{
		"set NeverCache, but it will be ignored":  manioc.NeverCache,
		"set ScopedCache, but it will be ignored": manioc.ScopedCache,
		"the default policy for RegisterInstance": manioc.GlobalCache,
	}
	for name, policy := range policies {
		policy := policy
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			ctr := manioc.NewContainer()
			err := manioc.RegisterInstance[IMyService](
				&MyService{Value: 42},
				manioc.WithContainer(ctr),
				manioc.WithCachePolicy(policy),
			)
			assert.Nil(err)

			ret1, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
			assert.Nil(err)
			assert.NotNil(ret1)
			ret2, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
			assert.Nil(err)
			assert.NotNil(ret2)

			// If you use RegisterInstance, an instance is created only once and will be reused.
			assert.Same(ret1, ret2)
		})
	}
}

func Test_CachePolicy_Validation(t *testing.T) {
	// Since the type CachePolicy is a named type of int, you can pass any int values as CachePolicy.
	// We treat this like an enum, so our implementation performs validation on runtime,
	// and causes panic for invalid CachePolicy values.
	// Once Golang introduces enum types, this could be checked statically.
	assert.Panics(t, func() {
		ctr := manioc.NewContainer()
		_ = manioc.Register[IMyService, MyService](
			manioc.WithContainer(ctr),
			manioc.WithCachePolicy(42 /* INVALID CachePolicy VALUE */),
		)
	})
}
