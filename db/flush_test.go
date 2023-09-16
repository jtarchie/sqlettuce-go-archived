package db_test

import (
	"context"

	"github.com/jtarchie/sqlettus/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flush", func() {
	var client *db.Client

	BeforeEach(func() {
		var err error

		client, err = db.NewClient("sqlite://:memory:?cache=shared&mode=memory")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		client.Close()
	})

	When("FlushAll", func() {
		It("resets all values", func() {
			err := client.Set(context.TODO(), "key", "value")
			Expect(err).NotTo(HaveOccurred())

			err = client.FlushAll(context.TODO())
			Expect(err).NotTo(HaveOccurred())

			value, found, err := client.Get(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
			Expect(value).To(Equal(""))
		})
	})
})
