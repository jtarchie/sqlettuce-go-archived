package db_test

import (
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

	When("Set", func() {
		It("sets a value", func() {
			err := client.Set("key", "value")
			Expect(err).NotTo(HaveOccurred())

			value, err := client.Get("key")
			Expect(err).NotTo(HaveOccurred())
			Expect(*value).To(Equal("value"))
		})
	})

	When("FlushAll", func() {
		It("resets all values", func() {
			err := client.Set("key", "value")
			Expect(err).NotTo(HaveOccurred())

			err = client.FlushAll()
			Expect(err).NotTo(HaveOccurred())

			value, err := client.Get("key")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(BeNil())
		})
	})
})
