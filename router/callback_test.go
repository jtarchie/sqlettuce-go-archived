package router_test

import (
	"bytes"
	"io"

	"github.com/jtarchie/sqlettuce/router"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Callback", func() {
	It("always returns the callback", func() {
		routes := router.CallbackRouter(func(tokens []string, writer io.Writer) error {
			Expect(tokens).To(Equal([]string{"HELLO"}))

			_, err := writer.Write([]byte("Hello"))
			Expect(err).NotTo(HaveOccurred())

			return nil
		})

		tokens := []string{"HELLO"}
		writer := &bytes.Buffer{}

		callback, found := routes.Lookup(tokens)
		Expect(found).To(BeTrue())

		err := callback(tokens, writer)
		Expect(err).NotTo(HaveOccurred())
		Expect(writer.String()).To(ContainSubstring("Hello"))
	})
})
