package router_test

import (
	"bytes"

	"github.com/jtarchie/sqlettus/router"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StaticResponse", func() {
	It("always writes a static response", func() {
		routes := router.StaticResponseRouter("Hello, World")

		callback, found := routes.Lookup(nil)
		Expect(found).To(BeTrue())

		writer := &bytes.Buffer{}

		err := callback(nil, writer)
		Expect(err).NotTo(HaveOccurred())
		Expect(writer.String()).To(ContainSubstring("Hello, World"))
	})
})
