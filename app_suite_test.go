package torm_test

import (
	"testing"

	"github.com/jackc/pgx/pgtype"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestModel struct {
	ID            pgtype.UUID `torm:"primary_key"`
	Foo           int64       `torm:"create_timestamp"`
	Bar           int64       `torm:"update_timestamp"`
	MyField       string
	Another_Field string
	Foo0          int
}

func TestScope(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Torm Test Suite")
}
