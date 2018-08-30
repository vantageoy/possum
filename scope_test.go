package torm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vantageoy/torm"
)

var _ = Describe("Scope", func() {

	var scope = torm.NewScope(&TestModel{})

	Describe("Scope Model", func() {

		It("should set primary field", func() {
			Expect(scope.Model.PrimaryField.Name).To(Equal("ID"))
		})

		It("should have a create timestamp field", func() {
			for _, field := range scope.Model.Fields {

				if field.HasCreateTimestampTag() {
					Expect(field.Name).To(Equal("Foo"))
					break
				}
			}
		})

		It("should have a update timestamp field", func() {
			for _, field := range scope.Model.Fields {

				if field.HasUpdateTimestampTag() {
					Expect(field.Name).To(Equal("Bar"))
					break
				}
			}
		})

		It("should return struct name", func() {
			Expect(scope.StructName()).To(Equal("TestModel"))
		})

		It("should return struct name in snake case", func() {
			Expect(scope.GetTableName()).To(Equal("test_models"))
		})

	})

	Describe("Scope sql", func() {

		It("should generate insert sql", func() {
			Expect(scope.CreateSQL()).To(Equal("insert into test_models (foo,bar,my_field,another__field,foo_0) VALUES ($1,$2,$3,$4,$5) RETURNING id"))
		})

	})

	Describe("Scope args", func() {

		It("should generate query args for insert", func() {

			model := TestModel{
				MyField:       "Foo Bar Biz b00",
				Another_Field: "Bar Biz Foo",
				Foo0:          11,
			}

			s := torm.NewScope(&model)

			// Expect timestamps to be 0, they are generated on insert
			Expect(s.CreateArgs()[0]).To(Equal(int64(0)))
			Expect(s.CreateArgs()[1]).To(Equal(int64(0)))

			Expect(s.CreateArgs()[2]).To(Equal(model.MyField))
			Expect(s.CreateArgs()[3]).To(Equal(model.Another_Field))
			Expect(s.CreateArgs()[4]).To(Equal(int64(11)))

		})

	})

})
