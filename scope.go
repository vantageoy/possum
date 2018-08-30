package torm

import (
	"fmt"
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

type QueryArgs []interface{}

var (
	CreateFieldFilter = FieldFilter{ExcludePrimary: true}
	UpdateFieldFilter = FieldFilter{ExcludeCreateTimestamp: true}
)

func NewScope(out interface{}) *Scope {

	scope := &Scope{Value: out}

	scope.Model = scope.GetModel()

	return scope

}

func (s *Scope) CreateSQL() string {

	unix := time.Now().Unix()

	if err := s.Model.SetCreateTimestamp(unix); err != nil {
		panic(err)
	}

	if err := s.Model.SetUpdateTimestamp(unix); err != nil {
		panic(err)
	}

	var colIndexes []string
	columns := s.Model.Columns(CreateFieldFilter)

	for i := range columns {
		colIndexes = append(colIndexes, fmt.Sprintf("$%s", strconv.Itoa(i+1)))
	}

	columnString := fmt.Sprintf("(%s)", strings.Join(columns, ","))
	valueString := fmt.Sprintf("VALUES (%s)", strings.Join(colIndexes, ","))

	fmt.Sprintf(columnString)
	fmt.Sprintf(valueString)

	return fmt.Sprintf("insert into %s %s %s RETURNING id", s.Model.GetTableName(s.Value), columnString, valueString)

}

func (s *Scope) CreateArgs() pgx.QueryArgs {

	return s.Model.Values(CreateFieldFilter)

}

func (s *Scope) FindSQL(where string) string {

	//columns := fmt.Sprintf("select %s from %s")

	return ""
}

func (s *Scope) CreateValuesSQL() string {

	var inputs []string

	for index, _ := range s.Model.Fields {
		inputs = append(inputs, fmt.Sprintf("$%d", index+1))
	}

	return fmt.Sprintf("values(%s)", strings.Join(inputs, ","))

}
