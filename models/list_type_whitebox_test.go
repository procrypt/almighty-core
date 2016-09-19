package models

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/almighty/almighty-core/resource"
	"github.com/stretchr/testify/assert"
)

func TestConvertList(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)

	dummyConverter := func(FieldType, interface{}) (interface{}, error) {
		return nil, nil
	}

	// Test case that value is not not an array or a string
	val1 := "foo"
	res, err := convertList(dummyConverter, SimpleType{Kind: KindString}, val1)
	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Equal(t, fmt.Sprintf(stErrorNotArrayOrSlice, val1, reflect.TypeOf(val1).Name()), err.Error())
}
