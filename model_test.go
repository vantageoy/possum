package possum_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vantageoy/possum"
)

var _ = Describe("Model", func() {

	var model = possum.NewScope(&TestModel{}).Model

	Describe("Model timestamps", func() {

		It("should set create timestamp", func() {

			unix := time.Now().Unix()

			model.SetCreateTimestamp(unix)

			for _, field := range model.Fields {
				if field.HasCreateTimestampTag() {
					Expect(field.RealValue()).To(Equal(int64(unix)))
				}
			}

		})

		It("should set update timestamp", func() {

			unix := time.Now().Unix()

			model.SetUpdateTimestamp(unix)

			for _, field := range model.Fields {
				if field.HasUpdateTimestampTag() {
					Expect(field.RealValue()).To(Equal(int64(unix)))
				}
			}

		})

	})

})
