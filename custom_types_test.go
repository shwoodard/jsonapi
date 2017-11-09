package jsonapi

import (
	"reflect"
	"testing"
)

func TestRegisterCustomTypes(t *testing.T) {
	for _, uuidType := range []reflect.Type{reflect.TypeOf(UUID{}), reflect.TypeOf(&UUID{})} {
		// given
		ClearCustomTypes() // make sure no other registration interferes with this test
		// when
		RegisterCustomTypeFunc(uuidType, reflect.TypeOf(""),
			func(value interface{}) (interface{}, error) {
				return "", nil
			},
			func(value interface{}) (interface{}, error) {
				return nil, nil
			})
		// then
		if !IsRegisteredType(uuidType) {
			t.Fatalf("Expected `%v` to be registered but it was not", uuidType)
		}
	}
}
