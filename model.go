package possum

import (
	"go/ast"
	"reflect"

	"github.com/jackc/pgx"
)

var (
	ErrNoCreateTimestamp = "Model has no field with create_timestamp tag"
	ErrNoUpdateTimestamp = "Model has no field with update_timestamp tag"
)

type FieldFilter struct {
	ExcludePrimary         bool
	ExcludeCreateTimestamp bool
	ExcludeUpdateTimeStamp bool
}

type Model struct {
	PrimaryField *Field
	Fields       []*Field
	Type         reflect.Type
}

func (s *Scope) GetModel() Model {

	var model Model

	if s.Value == nil {
		return model
	}

	model.Type = GetType(s.Value)
	modelValue := reflect.ValueOf(s.Value).Elem()

	for i := 0; i < modelValue.NumField(); i++ {

		if fieldStruct := model.Type.Field(i); ast.IsExported(fieldStruct.Name) {

			field := &Field{
				Name:   fieldStruct.Name,
				DBName: ToSnakeCase(fieldStruct.Name),
				Type:   model.Type.Field(i),
				Value:  modelValue.Field(i),
				Tag:    fieldStruct.Tag.Get("possum"),
			}

			if field.HasPrimaryTag() {
				model.PrimaryField = field
			}

			model.Fields = append(model.Fields, field)

		}

	}

	return model

}

func GetType(model interface{}) reflect.Type {

	reflectType := reflect.ValueOf(model).Type()

	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	if reflectType.Kind() != reflect.Struct {
		panic("Model value was not a struct")
	}

	return reflectType

}

func (m *Model) Values(filter FieldFilter) pgx.QueryArgs {

	args := pgx.QueryArgs{}

	for _, field := range m.Fields {

		if filter.ExcludePrimary && field.HasPrimaryTag() {
			continue
		}

		if filter.ExcludeCreateTimestamp && field.HasCreateTimestampTag() {
			continue
		}

		if filter.ExcludeUpdateTimeStamp && field.HasUpdateTimestampTag() {
			continue
		}

		args = append(args, field.RealValue())
	}

	return args

}

func (m *Model) Columns(filter FieldFilter) (columns []string) {

	for _, field := range m.Fields {

		if filter.ExcludePrimary && field.HasPrimaryTag() {
			continue
		}

		if filter.ExcludeCreateTimestamp && field.HasCreateTimestampTag() {
			continue
		}

		if filter.ExcludeUpdateTimeStamp && field.HasUpdateTimestampTag() {
			continue
		}

		columns = append(columns, field.DBName)
	}

	return columns

}

func (m *Model) SetCreateTimestamp(timestamp int64) error {

	for _, field := range m.Fields {

		if field.HasCreateTimestampTag() {
			if err := field.Set(timestamp); err == nil {
				return nil
			} else {
				return err
			}
		}

	}

	return nil

}

func (m *Model) SetUpdateTimestamp(timestamp int64) error {

	for _, field := range m.Fields {

		if field.HasUpdateTimestampTag() {
			if err := field.Set(timestamp); err == nil {
				return nil
			} else {
				return err
			}
		}

	}

	return nil

}
