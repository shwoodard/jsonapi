package jsonapi

import (
	"errors"
	"reflect"
)

var (
	ErrIncorrectInType  = errors.New("Incorrect input type supplied for marshal custom type")
	ErrIncorrectOutType = errors.New("Incorrect output type supplied for unmarshal custom type")
	customTypes         = map[reflect.Type]customType{}
)

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

func (ct customType) marshal(in interface{}) (interface{}, error) {
	if reflect.TypeOf(in) != ct.inType {
		return nil, ErrIncorrectInType
	}

	out, err := ct.marshaller.Marshal(in)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(out) != ct.outType {
		return nil, ErrIncorrectOutType
	}

	return out, nil
}

func (ct customType) unmarshal(out interface{}) (interface{}, error) {
	if reflect.TypeOf(out) != ct.outType {
		return nil, ErrIncorrectOutType
	}

	in, err := ct.unmarshaller.Unmarshal(out)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(in) != ct.inType {
		return nil, ErrIncorrectInType
	}

	return in, nil
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

// RegisterCustomType registers a custom type for use in marshalling and unmarshalling
func RegisterCustomType(
	inType, outType reflect.Type,
	marshaller TypeMarshaller,
	unmarshaller TypeUnmarshaller) {
	customTypes[inType] = newCustomType(inType, outType, marshaller, unmarshaller)
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
