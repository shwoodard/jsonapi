package jsonapi

import (
	"fmt"
	"reflect"
)

var (
	customTypes = map[reflect.Type]customType{}
)

type ErrIncorrectInType struct {
	expectedInType, incorrectInType reflect.Type
}

func (inErr *ErrIncorrectInType) Error() string {
	return fmt.Sprintf("Expected in type, %#v, got, %#v",
		inErr.expectedInType, inErr.incorrectInType)
}

func newErrIncorrectInType(
	expected, actual reflect.Type) *ErrIncorrectInType {
	return &ErrIncorrectInType{
		expectedInType:  expected,
		incorrectInType: actual,
	}
}

type ErrIncorrectOutType struct {
	expectedOutType, incorrectOutType reflect.Type
}

func (outErr *ErrIncorrectOutType) Error() string {
	return fmt.Sprintf("Expected out type, %#v, got, %#v",
		outErr.expectedOutType, outErr.incorrectOutType)
}

func newErrIncorrectOutType(
	expected, actual reflect.Type) *ErrIncorrectOutType {
	return &ErrIncorrectOutType{
		expectedOutType:  expected,
		incorrectOutType: actual,
	}
}

type customType struct {
	inType, outType reflect.Type
	marshaller      TypeMarshaller
	unmarshaller    TypeUnmarshaller
}

func newCustomType(
	inType, outType reflect.Type,
	marshaller TypeMarshaller,
	unmarshaller TypeUnmarshaller) customType {

	return customType{
		inType:       inType,
		outType:      outType,
		marshaller:   marshaller,
		unmarshaller: unmarshaller,
	}
}

func (ct customType) marshal(out interface{}) (interface{}, error) {
	inputType := reflect.TypeOf(out)
	if inputType != ct.outType {
		return nil, newErrIncorrectOutType(ct.outType, inputType)
	}

	in, err := ct.marshaller.Marshal(out)
	if err != nil {
		return nil, err
	}

	outputType := reflect.TypeOf(in)
	if outputType != ct.inType {
		return nil, newErrIncorrectInType(ct.inType, outputType)
	}

	return in, nil
}

func (ct customType) unmarshal(in interface{}) (interface{}, error) {
	inputType := reflect.TypeOf(in)
	if inputType != ct.inType {
		return nil, newErrIncorrectInType(ct.inType, inputType)
	}

	out, err := ct.unmarshaller.Unmarshal(in)
	if err != nil {
		return nil, err
	}

	outputType := reflect.TypeOf(out)
	if outputType != ct.outType {
		return nil, newErrIncorrectOutType(ct.outType, outputType)
	}

	return out, nil
}

type TypeMarshaller interface {
	Marshal(interface{}) (interface{}, error)
}

type TypeUnmarshaller interface {
	Unmarshal(interface{}) (interface{}, error)
}

type TypeMarshalFunc func(interface{}) (interface{}, error)

func (tmf TypeMarshalFunc) Marshal(in interface{}) (interface{}, error) {
	return tmf(in)
}

type TypeUnmarshalFunc func(interface{}) (interface{}, error)

func (tuf TypeUnmarshalFunc) Unmarshal(out interface{}) (interface{}, error) {
	return tuf(out)
}

// IsRegisteredType checks if the given type `t` is registered as a custom type
func IsRegisteredType(t reflect.Type) bool {
	_, ok := customTypes[t]
	return ok
}

// RegisterCustomType registers a custom type for use in marshalling and
// unmarshalling
func RegisterCustomType(
	inType, outType reflect.Type,
	marshaller TypeMarshaller,
	unmarshaller TypeUnmarshaller) {
	customTypes[inType] = newCustomType(
		inType, outType, marshaller, unmarshaller)
}

func RegisterCustomTypeFunc(
	inType, outType reflect.Type,
	marshaller TypeMarshalFunc,
	unmarshaller TypeUnmarshalFunc) {
	RegisterCustomType(inType, outType, marshaller, unmarshaller)
}

// ClearCustomTypes resets the custom type registration
func ClearCustomTypes() {
	customTypes = map[reflect.Type]customType{}
}
