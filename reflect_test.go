package manioc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type IFooService interface {
	DoFoo()
}

type IBarService interface {
	DoBar()
}

type IBazService interface {
	IFooService
	IBarService
	DoBaz()
}

type FooService struct{}

func (*FooService) DoFoo() {}

type BarService struct{}

func (*BarService) DoBar() {}

type BazService struct{}

func (*BazService) DoFoo() {}
func (*BazService) DoBar() {}
func (*BazService) DoBaz() {}

type BazService2 struct {
	FooService
	BarService
}

func (*BazService2) DoBaz() {}

func Test_ensureInterface_primitives(t *testing.T) {
	assert := assert.New(t)

	//
	// primitives
	//

	// boolean
	assert.Panics(func() { ensureInterface[bool]() })
	assert.Panics(func() { ensureInterface[*bool]() })
	// unsigned integer
	assert.Panics(func() { ensureInterface[uint]() })
	assert.Panics(func() { ensureInterface[*uint]() })
	assert.Panics(func() { ensureInterface[byte]() })
	assert.Panics(func() { ensureInterface[*byte]() })
	assert.Panics(func() { ensureInterface[uint8]() })
	assert.Panics(func() { ensureInterface[*uint8]() })
	assert.Panics(func() { ensureInterface[uint16]() })
	assert.Panics(func() { ensureInterface[*uint16]() })
	assert.Panics(func() { ensureInterface[uint32]() })
	assert.Panics(func() { ensureInterface[*uint32]() })
	assert.Panics(func() { ensureInterface[uint64]() })
	assert.Panics(func() { ensureInterface[*uint64]() })
	// signed integer
	assert.Panics(func() { ensureInterface[int]() })
	assert.Panics(func() { ensureInterface[*int]() })
	assert.Panics(func() { ensureInterface[int8]() })
	assert.Panics(func() { ensureInterface[*int8]() })
	assert.Panics(func() { ensureInterface[int16]() })
	assert.Panics(func() { ensureInterface[*int16]() })
	assert.Panics(func() { ensureInterface[int32]() })
	assert.Panics(func() { ensureInterface[*rune]() })
	assert.Panics(func() { ensureInterface[rune]() })
	assert.Panics(func() { ensureInterface[*int32]() })
	assert.Panics(func() { ensureInterface[int64]() })
	assert.Panics(func() { ensureInterface[*int64]() })
	// float
	assert.Panics(func() { ensureInterface[float32]() })
	assert.Panics(func() { ensureInterface[*float32]() })
	assert.Panics(func() { ensureInterface[float64]() })
	assert.Panics(func() { ensureInterface[*float64]() })
	// complex
	assert.Panics(func() { ensureInterface[complex64]() })
	assert.Panics(func() { ensureInterface[*complex64]() })
	assert.Panics(func() { ensureInterface[complex128]() })
	assert.Panics(func() { ensureInterface[*complex128]() })
	// string
	assert.Panics(func() { ensureInterface[string]() })
	assert.Panics(func() { ensureInterface[*string]() })
	// pointer
	assert.Panics(func() { ensureInterface[uintptr]() })
	assert.Panics(func() { ensureInterface[*uintptr]() })
}

func Test_ensureInterface_complex_types(t *testing.T) {
	assert := assert.New(t)

	// struct
	assert.Panics(func() { ensureInterface[struct{}]() })
	assert.Panics(func() { ensureInterface[*struct{}]() })
	assert.Panics(func() { ensureInterface[struct{ value int }]() })
	assert.Panics(func() { ensureInterface[*struct{ value int }]() })
	// array
	assert.Panics(func() { ensureInterface[[2]int]() })
	assert.Panics(func() { ensureInterface[*[2]int]() })
	assert.Panics(func() { ensureInterface[[2]struct{}]() })
	assert.Panics(func() { ensureInterface[*[2]struct{}]() })
	assert.Panics(func() { ensureInterface[[2]any]() })
	assert.Panics(func() { ensureInterface[*[2]any]() })
	// slice
	assert.Panics(func() { ensureInterface[[]int]() })
	assert.Panics(func() { ensureInterface[*[]int]() })
	assert.Panics(func() { ensureInterface[[]struct{}]() })
	assert.Panics(func() { ensureInterface[*[]struct{}]() })
	assert.Panics(func() { ensureInterface[[]any]() })
	assert.Panics(func() { ensureInterface[*[]any]() })
	// map
	assert.Panics(func() { ensureInterface[map[int]int]() })
	assert.Panics(func() { ensureInterface[*map[int]int]() })
	assert.Panics(func() { ensureInterface[map[struct{}]struct{}]() })
	assert.Panics(func() { ensureInterface[*map[struct{}]struct{}]() })
	assert.Panics(func() { ensureInterface[map[any]any]() })
	assert.Panics(func() { ensureInterface[*map[any]any]() })
	// function
	assert.Panics(func() { ensureInterface[func()]() })
	assert.Panics(func() { ensureInterface[*func()]() })
	assert.Panics(func() { ensureInterface[func() int]() })
	assert.Panics(func() { ensureInterface[*func() int]() })
	assert.Panics(func() { ensureInterface[func(int) int]() })
	assert.Panics(func() { ensureInterface[*func(int) int]() })
	// channel
	assert.Panics(func() { ensureInterface[chan int]() })
	assert.Panics(func() { ensureInterface[*chan int]() })
	assert.Panics(func() { ensureInterface[chan struct{}]() })
	assert.Panics(func() { ensureInterface[*chan struct{}]() })
	assert.Panics(func() { ensureInterface[chan any]() })
	assert.Panics(func() { ensureInterface[*chan any]() })
}

func Test_ensureInterface_interfaces(t *testing.T) {
	assert := assert.New(t)

	// interfaces
	assert.NotPanics(func() { ensureInterface[any]() })
	assert.NotPanics(func() { ensureInterface[any]() })
	assert.Panics(func() { ensureInterface[*any]() })
	assert.NotPanics(func() { ensureInterface[interface{ do() }]() })
	assert.Panics(func() { ensureInterface[*interface{ do() }]() })
	assert.NotPanics(func() { ensureInterface[error]() })

	// custom
	assert.NotPanics(func() { ensureInterface[IFooService]() })
	assert.Panics(func() { ensureInterface[*IFooService]() })
	assert.Panics(func() { ensureInterface[FooService]() })
	assert.Panics(func() { ensureInterface[*FooService]() })
}

func Test_ensureImplements(t *testing.T) {
	assert := assert.New(t)

	// arbitrary type implements any (any)
	assert.NotPanics(func() { ensureImplements[any, int]() })
	assert.NotPanics(func() { ensureImplements[any, struct{}]() })
	assert.NotPanics(func() { ensureImplements[any, *int]() })
	assert.NotPanics(func() { ensureImplements[any, *struct{}]() })
	assert.NotPanics(func() { ensureImplements[any, *any]() })
	assert.NotPanics(func() { ensureImplements[any, FooService]() })
	assert.NotPanics(func() { ensureImplements[any, BarService]() })
	assert.NotPanics(func() { ensureImplements[any, BazService]() })
	assert.NotPanics(func() { ensureImplements[any, BazService2]() })

	// do not pass non-interface type as TInterface
	assert.Panics(func() { ensureImplements[int, any]() })
	assert.Panics(func() { ensureImplements[int, int]() })
	assert.Panics(func() { ensureImplements[struct{}, any]() })
	assert.Panics(func() { ensureImplements[struct{}, struct{}]() })

	// do not pass interface type as TImplementation
	assert.Panics(func() { ensureImplements[any, any]() })
	assert.Panics(func() { ensureImplements[any, IFooService]() })
	assert.Panics(func() { ensureImplements[any, IBarService]() })
	assert.Panics(func() { ensureImplements[any, IBazService]() })

	// custom type
	assert.NotPanics(func() { ensureImplements[IFooService, FooService]() })
	assert.Panics(func() { ensureImplements[IFooService, BarService]() })
	assert.NotPanics(func() { ensureImplements[IFooService, BazService]() })
	assert.NotPanics(func() { ensureImplements[IFooService, BazService2]() })

	assert.Panics(func() { ensureImplements[IBarService, FooService]() })
	assert.NotPanics(func() { ensureImplements[IBarService, BarService]() })
	assert.NotPanics(func() { ensureImplements[IBarService, BazService]() })
	assert.NotPanics(func() { ensureImplements[IBarService, BazService2]() })

	assert.Panics(func() { ensureImplements[IBazService, FooService]() })
	assert.Panics(func() { ensureImplements[IBazService, BarService]() })
	assert.NotPanics(func() { ensureImplements[IBazService, BazService]() })
	assert.NotPanics(func() { ensureImplements[IBazService, BazService2]() })
}

func Test_ensureFunctionReturnType(t *testing.T) {
	assert := assert.New(t)

	// non function type
	assert.Panics(func() { ensureFunctionReturnType[any, any]() })
	assert.Panics(func() { ensureFunctionReturnType[int, any]() })
	assert.Panics(func() { ensureFunctionReturnType[*int, any]() })
	assert.Panics(func() { ensureFunctionReturnType[struct{}, any]() })

	// returns any
	assert.Panics(func() { ensureFunctionReturnType[func(), any]() })
	assert.NotPanics(func() { ensureFunctionReturnType[func() int, any]() })

	// args
	assert.Panics(func() { ensureFunctionReturnType[func(value int), any]() })
	assert.NotPanics(func() { ensureFunctionReturnType[func(value int) int, any]() })
	assert.NotPanics(func() { ensureFunctionReturnType[func(value int, value2 any) int, any]() })

	// varadic function is not allowed
	assert.Panics(func() { ensureFunctionReturnType[func(...int) any, any]() })

	// custom type
	assert.Panics(func() { ensureFunctionReturnType[func() FooService, IFooService]() })
	assert.NotPanics(func() { ensureFunctionReturnType[func() *FooService, IFooService]() })
	assert.Panics(func() { ensureFunctionReturnType[func() *FooService, IBarService]() })
	assert.Panics(func() { ensureFunctionReturnType[func() *FooService, IBazService]() })

	assert.Panics(func() { ensureFunctionReturnType[func() *BarService, IFooService]() })
	assert.NotPanics(func() { ensureFunctionReturnType[func() *BarService, IBarService]() })
	assert.Panics(func() { ensureFunctionReturnType[func() *BarService, IBazService]() })

	assert.NotPanics(func() { ensureFunctionReturnType[func() *BazService, IFooService]() })
	assert.NotPanics(func() { ensureFunctionReturnType[func() *BazService, IBarService]() })
	assert.NotPanics(func() { ensureFunctionReturnType[func() *BazService, IBazService]() })

	assert.NotPanics(func() { ensureFunctionReturnType[func() *BazService2, IFooService]() })
	assert.NotPanics(func() { ensureFunctionReturnType[func() *BazService2, IBarService]() })
	assert.NotPanics(func() { ensureFunctionReturnType[func() *BazService2, IBazService]() })
}
