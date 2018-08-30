package torm

import (
	"errors"
	"reflect"
	"strings"
)

var (
	// PrimaryTagName specifies if the field should be treated as primary key for its model.
	PrimaryTagName = "primary_key"
	// CreateTimestampTagName specifies if the field should be used to store creation time for parent model.
	CreateTimestampTagName = "create_timestamp"
	// UpdateTimestampTagName specifies if the field should be used to store timestamp when the parent model is updated.
	UpdateTimestampTagName = "update_timestamp"
)

type Field struct {
	DBName string
	Name   string
	Type   reflect.StructField
	Value  reflect.Value
	Tag    string
}

func (f *Field) Set(value interface{}) error {

	if !f.Value.IsValid() {
		return errors.New("field value is invalid")
	}

	if !f.Value.CanAddr() {
		return errors.New("field value is unaddressable")
	}

	reflectValue := reflect.ValueOf(value)

	fieldValue := f.Value

	if reflectValue.IsValid() {

		if reflectValue.Type().ConvertibleTo(fieldValue.Type()) {
			fieldValue.Set(reflectValue.Convert(fieldValue.Type()))
		}

	} else {
		f.Value.Set(reflect.Zero(f.Value.Type()))
	}

	return nil

}

func (f *Field) RealValue() interface{} {

	var value reflect.Value

	if f.Value.Kind() == reflect.Ptr {
		value = f.Value.Elem()
	} else {
		value = f.Value
	}

	switch value.Kind() {

	case reflect.String:
		return value.String()

	case reflect.Int:
		return value.Int()

	case reflect.Bool:
		return value.Bool()

	case reflect.Int16:
		return value.Int()

	case reflect.Int32:
		return value.Int()

	case reflect.Int64:
		return value.Int()

	}

	return nil

}

func (f *Field) HasTag(tag string) bool {

	if strings.Contains(f.Tag, tag) {
		return true
	}

	return false

}

func (f *Field) HasPrimaryTag() bool {

	return f.HasTag(PrimaryTagName)

}

func (f *Field) HasUpdateTimestampTag() bool {

	return f.HasTag(UpdateTimestampTagName)

}

func (f *Field) HasCreateTimestampTag() bool {

	return f.HasTag(CreateTimestampTagName)

}
