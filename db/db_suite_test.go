package db_test

import (
	"context"
	"testing"

	"github.com/jtarchie/sqlettus/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var _ = Describe("Client", func() {
	var client *db.Client

	BeforeEach(func() {
		var err error

		client, err = db.NewClient("file:test.db?cache=shared&mode=memory")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		client.Close()
	})

	When("Set", func() {
		It("sets a value", func() {
			err := client.Set(context.TODO(), "key", "value")
			Expect(err).NotTo(HaveOccurred())

			value, err := client.Get(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
			Expect(*value).To(Equal("value"))
		})
	})

	When("Get", func() {
		It("return nil with non existent keys", func() {
			value, err := client.Get(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(BeNil())
		})
	})

	When("Delete", func() {
		It("can delete a value", func() {
			err := client.Set(context.TODO(), "key", "value")
			Expect(err).NotTo(HaveOccurred())

			err = client.Delete(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())

			value, err := client.Get(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(BeNil())
		})

		It("does not fail on missing value", func() {
			err := client.Delete(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("FlushAll", func() {
		It("resets all values", func() {
			err := client.Set(context.TODO(), "key", "value")
			Expect(err).NotTo(HaveOccurred())

			err = client.FlushAll()
			Expect(err).NotTo(HaveOccurred())

			value, err := client.Get(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(BeNil())
		})
	})

	When("Append", func() {
		It("appends values to a string", func() {
			length, err := client.Append(context.TODO(), "name", "Hello")
			Expect(err).NotTo(HaveOccurred())
			Expect(length).To(BeEquivalentTo(5))

			value, err := client.Get(context.TODO(), "name")
			Expect(err).NotTo(HaveOccurred())
			Expect(*value).To(Equal("Hello"))

			length, err = client.Append(context.TODO(), "name", " World")
			Expect(err).NotTo(HaveOccurred())
			Expect(length).To(BeEquivalentTo(11))

			value, err = client.Get(context.TODO(), "name")
			Expect(err).NotTo(HaveOccurred())
			Expect(*value).To(Equal("Hello World"))
		})
	})
})
