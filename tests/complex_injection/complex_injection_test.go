package manioc_complex_injection_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/fuzmish/manioc"
	"github.com/stretchr/testify/assert"
)

// In this case, we will verify that we can resolve services
// with complex dependencies such as the following:
//  A implements IA
//  B implements IB
//    B requires IA via field injection
//  C implements IC
//    C requires IB via constructor injection
//  D implements ID
//    D requires IA via constructor injection
//    D requires IB via field injection
//  E implements IE
//   E requires []IA via field injection
//   E requires []IB via constructor injection
//   E requires IC via field injection
//   E requires ID via constructor injection

type IA interface {
	doA()
}

type A struct{}

func (a *A) doA() {}

type IB interface {
	doB()
}

type B struct {
	a IA `manioc:"inject"`
}

func (b *B) doB() {}

type IC interface {
	doC()
}

type C struct {
	b IB
}

func (c *C) doC() {}

func NewC(b IB) *C {
	return &C{b: b}
}

type ID interface {
	doD()
}

type D struct {
	a IA
	b IB `manioc:"inject"`
}

func (d *D) doD() {}

func NewD(a IA) *D {
	//nolint
	return &D{a: a}
}

type IE interface {
	doE()
}

type E struct {
	a []IA `manioc:"inject"`
	b []IB
	c IC `manioc:"inject"`
	d ID
}

func (e *E) doE() {}

func NewE(b []IB, d ID) *E {
	//nolint
	return &E{b: b, d: d}
}

func Test_ComplexInjection_Errors(t *testing.T) {
	// test all possible registration patterns
	numDependencies := 5
	for i := 0; i < int(math.Pow(2, float64(numDependencies))); i++ {
		pattern := i
		t.Run(fmt.Sprintf("pattern: %d", i), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			// register
			ctr := manioc.NewContainer()
			numRegistered := 0
			for j := 0; j < numDependencies; j++ {
				if pattern>>j == 0 {
					continue
				}
				switch j {
				case 0:
					assert.Nil(manioc.RegisterSingleton[IA, A](manioc.WithContainer(ctr)))
				case 1:
					assert.Nil(manioc.RegisterSingleton[IB, B](manioc.WithContainer(ctr)))
				case 2:
					assert.Nil(manioc.RegisterSingletonConstructor[IC](NewC, manioc.WithContainer(ctr)))
				case 3:
					assert.Nil(manioc.RegisterSingletonConstructor[ID](NewD, manioc.WithContainer(ctr)))
				case 4:
					assert.Nil(manioc.RegisterConstructor[IE](NewE, manioc.WithContainer(ctr)))
				}
				numRegistered++
			}

			// resolve
			ret, err := manioc.Resolve[IE](manioc.WithScope(ctr))
			if numRegistered == numDependencies {
				assert.NotNil(ret)
				assert.Nil(err)
			} else {
				assert.Nil(ret)
				assert.Error(err)
			}
		})
	}
}

func Test_ComplexInjection(t *testing.T) {
	assert := assert.New(t)

	// register
	ctr := manioc.NewContainer()
	assert.Nil(manioc.RegisterSingleton[IA, A](manioc.WithContainer(ctr)))
	assert.Nil(manioc.RegisterSingleton[IB, B](manioc.WithContainer(ctr)))
	assert.Nil(manioc.RegisterSingletonConstructor[IC](NewC, manioc.WithContainer(ctr)))
	assert.Nil(manioc.RegisterSingletonConstructor[ID](NewD, manioc.WithContainer(ctr)))
	assert.Nil(manioc.RegisterConstructor[IE](NewE, manioc.WithContainer(ctr)))

	// resolve
	ret, err := manioc.Resolve[IE](manioc.WithScope(ctr))
	assert.NotNil(ret)
	assert.Nil(err)

	refA := manioc.MustResolve[IA](manioc.WithScope(ctr))
	refB := manioc.MustResolve[IB](manioc.WithScope(ctr))
	refC := manioc.MustResolve[IC](manioc.WithScope(ctr))
	refD := manioc.MustResolve[ID](manioc.WithScope(ctr))

	// check injected instances

	//nolint
	e, ok := ret.(*E)
	assert.True(ok)
	assert.Len(e.a, 1)
	assert.Same(e.a[0], refA)
	assert.Len(e.b, 1)
	assert.Same(e.b[0], refB)
	assert.NotNil(e.c)
	assert.Same(e.c, refC)
	assert.NotNil(e.d)
	assert.Same(e.d, refD)

	d, ok := e.d.(*D)
	assert.True(ok)
	assert.NotNil(d.a)
	assert.Same(d.a, refA)
	assert.NotNil(d.b)
	assert.Same(d.b, refB)

	c, ok := e.c.(*C)
	assert.True(ok)
	assert.NotNil(c.b)
	assert.Same(c.b, refB)

	b, ok := c.b.(*B)
	assert.True(ok)
	assert.NotNil(b.a)
	assert.Same(b.a, refA)
}
