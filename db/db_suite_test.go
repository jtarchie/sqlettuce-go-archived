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
	When("Set", func() {
		It("sets a value", func() {
			client, err := db.NewClient("file:test.db?cache=shared&mode=memory")
			Expect(err).NotTo(HaveOccurred())

			err = client.Set("key", "value")
			Expect(err).NotTo(HaveOccurred())

			value, err := client.Get("key")
			Expect(err).NotTo(HaveOccurred())
			Expect(*value).To(Equal("value"))
		})
	})
})
