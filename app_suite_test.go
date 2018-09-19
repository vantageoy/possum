package possum_test

import (
	"testing"

	"github.com/jackc/pgx/pgtype"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestModel struct {
	ID            pgtype.UUID `possum:"primary_key"`
	Foo           int64       `possum:"create_timestamp"`
	Bar           int64       `possum:"update_timestamp"`
	MyField       string
	Another_Field string
	Foo0          int
}

func TestScope(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "possum Test Suite")
}
