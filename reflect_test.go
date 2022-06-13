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

func Test_ensureInterface(t *testing.T) {
	assert := assert.New(t)

	//
	// primitives
	//

	// boolean
	assert.Error(ensureInterface[bool]())
	assert.Error(ensureInterface[*bool]())
	// unsigned integer
	assert.Error(ensureInterface[uint]())
	assert.Error(ensureInterface[*uint]())
	assert.Error(ensureInterface[byte]())
	assert.Error(ensureInterface[*byte]())
	assert.Error(ensureInterface[uint8]())
	assert.Error(ensureInterface[*uint8]())
	assert.Error(ensureInterface[uint16]())
	assert.Error(ensureInterface[*uint16]())
	assert.Error(ensureInterface[uint32]())
	assert.Error(ensureInterface[*uint32]())
	assert.Error(ensureInterface[uint64]())
	assert.Error(ensureInterface[*uint64]())
	// signed integer
	assert.Error(ensureInterface[int]())
	assert.Error(ensureInterface[*int]())
	assert.Error(ensureInterface[int8]())
	assert.Error(ensureInterface[*int8]())
	assert.Error(ensureInterface[int16]())
	assert.Error(ensureInterface[*int16]())
	assert.Error(ensureInterface[int32]())
	assert.Error(ensureInterface[*rune]())
	assert.Error(ensureInterface[rune]())
	assert.Error(ensureInterface[*int32]())
	assert.Error(ensureInterface[int64]())
	assert.Error(ensureInterface[*int64]())
	// float
	assert.Error(ensureInterface[float32]())
	assert.Error(ensureInterface[*float32]())
	assert.Error(ensureInterface[float64]())
	assert.Error(ensureInterface[*float64]())
	// complex
	assert.Error(ensureInterface[complex64]())
	assert.Error(ensureInterface[*complex64]())
	assert.Error(ensureInterface[complex128]())
	assert.Error(ensureInterface[*complex128]())
	// string
	assert.Error(ensureInterface[string]())
	assert.Error(ensureInterface[*string]())
	// pointer
	assert.Error(ensureInterface[uintptr]())
	assert.Error(ensureInterface[*uintptr]())

	//
	// interfaces
	//
	assert.NoError(ensureInterface[any]())
	assert.NoError(ensureInterface[any]())
	assert.Error(ensureInterface[*any]())
	assert.NoError(ensureInterface[interface{ do() }]())
	assert.Error(ensureInterface[*interface{ do() }]())
	assert.NoError(ensureInterface[error]())

	//
	// complex types
	//
	// struct
	assert.Error(ensureInterface[struct{}]())
	assert.Error(ensureInterface[*struct{}]())
	assert.Error(ensureInterface[struct{ value int }]())
	assert.Error(ensureInterface[*struct{ value int }]())
	// array
	assert.Error(ensureInterface[[2]int]())
	assert.Error(ensureInterface[*[2]int]())
	assert.Error(ensureInterface[[2]struct{}]())
	assert.Error(ensureInterface[*[2]struct{}]())
	assert.Error(ensureInterface[[2]any]())
	assert.Error(ensureInterface[*[2]any]())
	// slice
	assert.Error(ensureInterface[[]int]())
	assert.Error(ensureInterface[*[]int]())
	assert.Error(ensureInterface[[]struct{}]())
	assert.Error(ensureInterface[*[]struct{}]())
	assert.Error(ensureInterface[[]any]())
	assert.Error(ensureInterface[*[]any]())
	// map
	assert.Error(ensureInterface[map[int]int]())
	assert.Error(ensureInterface[*map[int]int]())
	assert.Error(ensureInterface[map[struct{}]struct{}]())
	assert.Error(ensureInterface[*map[struct{}]struct{}]())
	assert.Error(ensureInterface[map[any]any]())
	assert.Error(ensureInterface[*map[any]any]())
	// function
	assert.Error(ensureInterface[func()]())
	assert.Error(ensureInterface[*func()]())
	assert.Error(ensureInterface[func() int]())
	assert.Error(ensureInterface[*func() int]())
	assert.Error(ensureInterface[func(int) int]())
	assert.Error(ensureInterface[*func(int) int]())
	// channel
	assert.Error(ensureInterface[chan int]())
	assert.Error(ensureInterface[*chan int]())
	assert.Error(ensureInterface[chan struct{}]())
	assert.Error(ensureInterface[*chan struct{}]())
	assert.Error(ensureInterface[chan any]())
	assert.Error(ensureInterface[*chan any]())

	// custom
	assert.NoError(ensureInterface[IFooService]())
	assert.Error(ensureInterface[*IFooService]())
	assert.Error(ensureInterface[FooService]())
	assert.Error(ensureInterface[*FooService]())
}

func Test_ensureImplements(t *testing.T) {
	assert := assert.New(t)

	// arbitrary type implements any (any)
	assert.NoError(ensureImplements[any, int]())
	assert.NoError(ensureImplements[any, struct{}]())
	assert.NoError(ensureImplements[any, *int]())
	assert.NoError(ensureImplements[any, *struct{}]())
	assert.NoError(ensureImplements[any, *any]())
	assert.NoError(ensureImplements[any, FooService]())
	assert.NoError(ensureImplements[any, BarService]())
	assert.NoError(ensureImplements[any, BazService]())
	assert.NoError(ensureImplements[any, BazService2]())

	// do not pass non-interface type as TInterface
	assert.Error(ensureImplements[int, any]())
	assert.Error(ensureImplements[int, int]())
	assert.Error(ensureImplements[struct{}, any]())
	assert.Error(ensureImplements[struct{}, struct{}]())

	// do not pass interface type as TImplementation
	assert.Error(ensureImplements[any, any]())
	assert.Error(ensureImplements[any, IFooService]())
	assert.Error(ensureImplements[any, IBarService]())
	assert.Error(ensureImplements[any, IBazService]())

	// custom type
	assert.NoError(ensureImplements[IFooService, FooService]())
	assert.Error(ensureImplements[IFooService, BarService]())
	assert.NoError(ensureImplements[IFooService, BazService]())
	assert.NoError(ensureImplements[IFooService, BazService2]())

	assert.Error(ensureImplements[IBarService, FooService]())
	assert.NoError(ensureImplements[IBarService, BarService]())
	assert.NoError(ensureImplements[IBarService, BazService]())
	assert.NoError(ensureImplements[IBarService, BazService2]())

	assert.Error(ensureImplements[IBazService, FooService]())
	assert.Error(ensureImplements[IBazService, BarService]())
	assert.NoError(ensureImplements[IBazService, BazService]())
	assert.NoError(ensureImplements[IBazService, BazService2]())
}

func Test_ensureFunctionReturnType(t *testing.T) {
	assert := assert.New(t)

	// non function type
	assert.Error(ensureFunctionReturnType[any, any]())
	assert.Error(ensureFunctionReturnType[int, any]())
	assert.Error(ensureFunctionReturnType[*int, any]())
	assert.Error(ensureFunctionReturnType[struct{}, any]())

	// returns any
	assert.Error(ensureFunctionReturnType[func(), any]())
	assert.NoError(ensureFunctionReturnType[func() int, any]())

	// args
	assert.Error(ensureFunctionReturnType[func(value int), any]())
	assert.NoError(ensureFunctionReturnType[func(value int) int, any]())
	assert.NoError(ensureFunctionReturnType[func(value int, value2 any) int, any]())

	// varadic function is not allowed
	assert.Error(ensureFunctionReturnType[func(...int) any, any]())

	// custom type
	assert.Error(ensureFunctionReturnType[func() FooService, IFooService]())
	assert.NoError(ensureFunctionReturnType[func() *FooService, IFooService]())
	assert.Error(ensureFunctionReturnType[func() *FooService, IBarService]())
	assert.Error(ensureFunctionReturnType[func() *FooService, IBazService]())

	assert.Error(ensureFunctionReturnType[func() *BarService, IFooService]())
	assert.NoError(ensureFunctionReturnType[func() *BarService, IBarService]())
	assert.Error(ensureFunctionReturnType[func() *BarService, IBazService]())

	assert.NoError(ensureFunctionReturnType[func() *BazService, IFooService]())
	assert.NoError(ensureFunctionReturnType[func() *BazService, IBarService]())
	assert.NoError(ensureFunctionReturnType[func() *BazService, IBazService]())

	assert.NoError(ensureFunctionReturnType[func() *BazService2, IFooService]())
	assert.NoError(ensureFunctionReturnType[func() *BazService2, IBarService]())
	assert.NoError(ensureFunctionReturnType[func() *BazService2, IBazService]())
}
