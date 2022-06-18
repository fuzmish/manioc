package manioc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertParseTagResult(t *testing.T, data any, fieldIndex int, expected *tagInfo) {
	t.Helper()
	ret, err := parseTag(reflect.TypeOf(data).Field(fieldIndex).Tag)
	assert.Nil(t, err)
	assert.Equal(t, expected, ret)
}

func assertParseTagError(t *testing.T, data any, fieldIndex int) {
	t.Helper()
	_, err := parseTag(reflect.TypeOf(data).Field(fieldIndex).Tag)
	assert.Error(t, err)
}

func Test_parseTag_Empty(t *testing.T) {
	t.Run("when the field has no tags, returns inject=false", func(t *testing.T) {
		var data struct {
			value0 any
			value1 any `manioc:""`
		}
		assertParseTagResult(t, data, 0, &tagInfo{inject: false, key: nil})
		assertParseTagResult(t, data, 1, &tagInfo{inject: false, key: nil})
	})
}

func Test_parseTag_Inject(t *testing.T) {
	t.Run("manioc: inject", func(t *testing.T) {
		var data struct {
			value any `manioc:"inject"`
		}
		assertParseTagResult(t, data, 0, &tagInfo{inject: true, key: nil})
	})

	t.Run("duplicated inject tag is allowed", func(t *testing.T) {
		var data struct {
			value any `manioc:"inject,inject"`
		}
		assertParseTagResult(t, data, 0, &tagInfo{inject: true, key: nil})
	})
}

func Test_parseTag_Key(t *testing.T) {
	t.Run("key requires value", func(t *testing.T) {
		var data struct {
			value any `manioc:"key"`
		}
		assertParseTagError(t, data, 0)
	})

	t.Run("key without value; an empyt string key is not allowed", func(t *testing.T) {
		var data struct {
			value any `manioc:"key="`
		}
		assertParseTagResult(t, data, 0, &tagInfo{inject: false, key: nil})
	})

	t.Run("key with value", func(t *testing.T) {
		var data struct {
			value any `manioc:"key=foo"`
		}
		assertParseTagResult(t, data, 0, &tagInfo{inject: false, key: "foo"})
	})

	t.Run("if key tag is duplicated, the last value is used without error", func(t *testing.T) {
		var data struct {
			value any `manioc:"key=foo,key=bar"`
		}
		assertParseTagResult(t, data, 0, &tagInfo{inject: false, key: "bar"})
	})
}

func Test_parseTag_Combinations(t *testing.T) {
	t.Run("valid tags", func(t *testing.T) {
		var data struct {
			value0 any `manioc:"inject,key=foo"`
			value1 any `manioc:"key=foo,inject"`
			value2 any `manioc:"inject,key=foo,inject"`
			value3 any `manioc:"key=foo,inject,key=bar"`
		}
		assertParseTagResult(t, data, 0, &tagInfo{inject: true, key: "foo"})
		assertParseTagResult(t, data, 1, &tagInfo{inject: true, key: "foo"})
		assertParseTagResult(t, data, 2, &tagInfo{inject: true, key: "foo"})
		assertParseTagResult(t, data, 3, &tagInfo{inject: true, key: "bar"})
	})

	t.Run("with other tags", func(t *testing.T) {
		var data struct {
			Value0 any `json:"value" manioc:"inject"`
			Value1 any `manioc:"key=foo" gorm:"index"`
		}
		assertParseTagResult(t, data, 0, &tagInfo{inject: true, key: nil})
		assertParseTagResult(t, data, 1, &tagInfo{inject: false, key: "foo"})
	})
}

func Test_parseTag_Validation(t *testing.T) {
	t.Run("returns error for unknown tags", func(t *testing.T) {
		var data struct {
			value any `manioc:"unknownTag"`
		}
		assertParseTagError(t, data, 0)
	})

	t.Run("trailing comma is allowed", func(t *testing.T) {
		var data struct {
			value0 any `manioc:"inject,"`
			value1 any `manioc:"inject,key=foo,"`
		}
		assertParseTagResult(t, data, 0, &tagInfo{inject: true, key: nil})
		assertParseTagResult(t, data, 1, &tagInfo{inject: true, key: "foo"})
	})
}
