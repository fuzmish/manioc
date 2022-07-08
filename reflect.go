package manioc

import (
	"reflect"
)

func typeof[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func nameof[T any]() string {
	t := typeof[T]()
	ret := t.PkgPath()
	if ret != "" {
		ret += "."
	}
	ret += t.Name()
	return ret
}
