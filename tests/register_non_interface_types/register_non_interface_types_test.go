package manioc_register_non_interface_types_test

import (
	"testing"

	"github.com/fuzmish/manioc"
	"github.com/stretchr/testify/assert"
)

func Test_Register_PrimitiveValue(t *testing.T) {
	t.Run("register primitive type", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[int, int](manioc.WithContainer(ctr)))
		assert.Zero(manioc.MustResolve[int](manioc.WithScope(ctr)))
	})

	t.Run("register primitive pointer type", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[*int, *int](manioc.WithContainer(ctr)))
		assert.Zero(*manioc.MustResolve[*int](manioc.WithScope(ctr)))
	})

	t.Run("register primitive constructor", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterConstructor[int](func() int { return 42 }, manioc.WithContainer(ctr)))
		assert.Equal(42, manioc.MustResolve[int](manioc.WithScope(ctr)))
	})

	t.Run("register primitive value", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterInstance(42, manioc.WithContainer(ctr)))
		assert.Equal(42, manioc.MustResolve[int](manioc.WithScope(ctr)))
	})

	t.Run("register primitive multiple value", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterInstance(42, manioc.WithContainer(ctr)))
		assert.Nil(manioc.RegisterInstance(84, manioc.WithContainer(ctr)))
		assert.Equal([]int{42, 84}, manioc.MustResolveMany[int](manioc.WithScope(ctr)))
	})
}

func Test_Register_ArrayLike(t *testing.T) {
	t.Run("register array", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterInstance([4]int{0, 1, 2, 3}, manioc.WithContainer(ctr)))
		assert.Equal([4]int{0, 1, 2, 3}, manioc.MustResolve[[4]int](manioc.WithScope(ctr)))
	})

	t.Run("register slice", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterInstance([]int{0, 1, 2, 3}, manioc.WithContainer(ctr)))
		assert.Equal([]int{0, 1, 2, 3}, manioc.MustResolve[[]int](manioc.WithScope(ctr)))
	})
}

func Test_Register_Map(t *testing.T) {
	t.Run("register map", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterInstance(map[int]string{42: "hello world"}, manioc.WithContainer(ctr)))
		assert.Equal(map[int]string{42: "hello world"}, manioc.MustResolve[map[int]string](manioc.WithScope(ctr)))
	})
}

func Test_Register_Func(t *testing.T) {
	t.Run("register func", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterInstance(func(s string) int { return len(s) }, manioc.WithContainer(ctr)))
		fn := manioc.MustResolve[func(string) int](manioc.WithScope(ctr))
		assert.NotNil(fn)
		assert.Equal(5, fn("hello"))
	})
}

type MyObject struct {
	property int
}

func Test_Register_Struct(t *testing.T) {
	t.Run("register struct type", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[MyObject, MyObject](manioc.WithContainer(ctr)))
		ret := manioc.MustResolve[MyObject](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Equal(0, ret.property)
	})

	t.Run("register struct pointer type", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.Register[*MyObject, *MyObject](manioc.WithContainer(ctr)))
		ret := manioc.MustResolve[*MyObject](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Equal(0, ret.property)
	})

	t.Run("register struct value", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterInstance(MyObject{property: 42}, manioc.WithContainer(ctr)))
		ret := manioc.MustResolve[MyObject](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Equal(42, ret.property)
	})

	t.Run("register struct pointer value", func(t *testing.T) {
		assert := assert.New(t)

		ctr := manioc.NewContainer()
		assert.Nil(manioc.RegisterInstance(&MyObject{property: 42}, manioc.WithContainer(ctr)))
		ret := manioc.MustResolve[*MyObject](manioc.WithScope(ctr))
		assert.NotNil(ret)
		assert.Equal(42, ret.property)
	})
}
