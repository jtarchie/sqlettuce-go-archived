package db_test

import (
	"context"

	"github.com/jtarchie/sqlettus/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integer", func() {
	var client *db.Client

	BeforeEach(func() {
		var err error

		client, err = db.NewClient(":memory:?cache=shared&mode=memory")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		client.Close()
	})

	It("can increment and decrement a number", func() {
		err := client.Set(context.Background(), "key", "10")
		Expect(err).NotTo(HaveOccurred())

		intValue, err := client.AddInt(context.Background(), "key", 1)
		Expect(err).NotTo(HaveOccurred())
		Expect(intValue).To(BeEquivalentTo(11))

		value, err := client.Get(context.Background(), "key")
		Expect(err).NotTo(HaveOccurred())
		Expect(*value).To(Equal("11"))

		intValue, err = client.AddInt(context.Background(), "key", -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(intValue).To(BeEquivalentTo(10))

		value, err = client.Get(context.Background(), "key")
		Expect(err).NotTo(HaveOccurred())
		Expect(*value).To(Equal("10"))
	})

	When("the value is non integer string", func() {
		It("returns an error", func() {
			err := client.Set(context.Background(), "key", "value")
			Expect(err).NotTo(HaveOccurred())

			_, err = client.AddInt(context.Background(), "key", 1)
			Expect(err).To(HaveOccurred())
		})
	})
})
