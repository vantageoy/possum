package possum

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx"
)

type Scope struct {
	Name  string
	Value interface{}
	Model Model
}

var (
	CreateFieldFilter = FieldFilter{ExcludePrimary: true}
	UpdateFieldFilter = FieldFilter{ExcludePrimary: true, ExcludeCreateTimestamp: true}
)

func NewScope(out interface{}) *Scope {

	scope := &Scope{Value: out}

	scope.Model = scope.GetModel()

	return scope

}

func (s *Scope) StructName() string {

	return s.GetStruct(s.Value).Name()

}

func (s *Scope) GetStruct(model interface{}) reflect.Type {

	modelStruct := reflect.ValueOf(model).Type()
	for modelStruct.Kind() == reflect.Slice || modelStruct.Kind() == reflect.Ptr {
		modelStruct = modelStruct.Elem()
	}

	// Scope value needs to be a struct
	if modelStruct.Kind() != reflect.Struct {
		panic("Model value was not a struct")
	}

	return modelStruct

}

func (s *Scope) GetTableName() string {

	return fmt.Sprintf("%ss", ToSnakeCase(s.StructName()))

}

func (s *Scope) CreateSQL() string {

	unix := timestamp()

	if err := s.Model.SetCreateTimestamp(unix); err != nil {
		panic(err)
	}

	if err := s.Model.SetUpdateTimestamp(unix); err != nil {
		panic(err)
	}

	columns := s.Model.Columns(CreateFieldFilter)
	colIndexes := s.ColumnIndexes(columns)

	columnString := fmt.Sprintf("(%s)", strings.Join(columns, ","))
	valueString := fmt.Sprintf("VALUES (%s)", strings.Join(colIndexes, ","))

	return fmt.Sprintf("insert into %s %s %s RETURNING id", s.GetTableName(), columnString, valueString)

}

func (s *Scope) CreateArgs() pgx.QueryArgs {

	return s.Model.Values(CreateFieldFilter)

}

func (s *Scope) ColumnIndexes(columns []string) (indexes []string) {

	var colIndexes []string

	for i := range columns {
		colIndexes = append(colIndexes, fmt.Sprintf("$%s", strconv.Itoa(i+1)))
	}

	return colIndexes

}

func timestamp() int64 {

	return time.Now().Unix()

}
