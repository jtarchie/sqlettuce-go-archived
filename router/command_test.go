package router_test

import (
	"bytes"
	"io"

	"github.com/jtarchie/sqlettuce/router"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Command", func() {
	When("using commands", func() {
		It("looks up a command and returns in callback", func() {
			routes := router.Command{
				"HELLO": router.CallbackRouter(func(tokens []string, writer io.Writer) error {
					Expect(tokens).To(Equal([]string{"HELLO"}))

					_, err := writer.Write([]byte("Hello"))
					Expect(err).NotTo(HaveOccurred())

					return nil
				}),
			}

			tokens := []string{"HELLO"}

			callback, found := routes.Lookup(tokens)
			Expect(found).To(BeTrue())

			writer := &bytes.Buffer{}

			err := callback(tokens, writer)
			Expect(err).NotTo(HaveOccurred())
			Expect(writer.String()).To(Equal("Hello"))
		})

		When("cannot find a command", func() {
			It("writes an unsupported command", func() {
				routes := router.Command{}
				tokens := []string{"HELLO"}
				writer := &bytes.Buffer{}

				callback, found := routes.Lookup(tokens)
				Expect(found).To(BeFalse())

				err := callback(tokens, writer)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(ContainSubstring("Unsupported command"))
			})
		})

		When("supporting nested commands", func() {
			It("keeps looking up till found", func() {
				routes := router.Command{
					"HELLO": router.Command{
						"WORLD": router.CallbackRouter(func(tokens []string, writer io.Writer) error {
							Expect(tokens).To(Equal([]string{"HELLO", "WORLD"}))

							_, err := writer.Write([]byte("Hello"))
							Expect(err).NotTo(HaveOccurred())

							return nil
						}),
					},
				}

				tokens := []string{"HELLO", "WORLD"}

				callback, found := routes.Lookup(tokens)
				Expect(found).To(BeTrue())

				writer := &bytes.Buffer{}

				err := callback(tokens, writer)
				Expect(err).NotTo(HaveOccurred())
				Expect(writer.String()).To(Equal("Hello"))
			})

			When("sub commands cannot be found", func() {
				It("writes an unsupported command", func() {
					routes := router.Command{
						"HELLO": router.Command{},
					}

					tokens := []string{"HELLO", "WORLD"}

					callback, found := routes.Lookup(tokens)
					Expect(found).To(BeFalse())

					writer := &bytes.Buffer{}

					err := callback(tokens, writer)
					Expect(err).NotTo(HaveOccurred())
					Expect(writer.String()).To(ContainSubstring("Unsupported command"))
				})
			})
		})
	})
})
