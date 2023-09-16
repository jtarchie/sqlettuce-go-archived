package db_test

import (
	"context"

	"github.com/jtarchie/sqlettus/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Float", func() {
	var client *db.Client

	BeforeEach(func() {
		var err error

		client, err = db.NewClient("sqlite://:memory:?cache=shared&mode=memory")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		client.Close()
	})

	It("can increment and decrement a number", func() {
		err := client.Set(context.Background(), "key", "10.50")
		Expect(err).NotTo(HaveOccurred())

		value, found, err := client.Get(context.Background(), "key")
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue())
		Expect(value).To(Equal("10.50"))

		floatValue, err := client.AddFloat(context.Background(), "key", 0.1)
		Expect(err).NotTo(HaveOccurred())
		Expect(floatValue).To(BeEquivalentTo(10.6))

		value, found, err = client.Get(context.Background(), "key")
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue())
		Expect(value).To(Equal("10.6"))

		floatValue, err = client.AddFloat(context.Background(), "key", -5)
		Expect(err).NotTo(HaveOccurred())
		Expect(floatValue).To(BeEquivalentTo(5.6))

		value, found, err = client.Get(context.Background(), "key")
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue())
		Expect(value).To(Equal("5.6"))

		err = client.Set(context.Background(), "key", "5.0e3")
		Expect(err).NotTo(HaveOccurred())

		floatValue, err = client.AddFloat(context.Background(), "key", 2.0e2)
		Expect(err).NotTo(HaveOccurred())
		Expect(floatValue).To(BeEquivalentTo(5200))
	})

	When("the value is non integer string", func() {
		It("returns an error", func() {
			err := client.Set(context.Background(), "key", "value")
			Expect(err).NotTo(HaveOccurred())

			_, err = client.AddFloat(context.Background(), "key", 1)
			Expect(err).To(HaveOccurred())
		})
	})
})
