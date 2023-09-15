package router_test

import (
	"bytes"
	"io"

	"github.com/jtarchie/sqlettus/router"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Limits", func() {
	When("setting a minimum", func() {
		It("requires that minimum or cannot be found", func() {
			routes := router.MinMaxTokensRouter(1, 0, func(s []string, w io.Writer) error {
				return nil
			})

			writer := &bytes.Buffer{}

			callback, found := routes.Lookup([]string{})
			Expect(found).To(BeFalse())

			err := callback(nil, writer)
			Expect(err).NotTo(HaveOccurred())
			Expect(writer.String()).To(ContainSubstring("expected minimum"))

			writer.Reset()

			callback, found = routes.Lookup([]string{"HELLO"})
			Expect(found).To(BeTrue())

			err = callback(nil, writer)
			Expect(err).NotTo(HaveOccurred())
			Expect(writer.String()).NotTo(ContainSubstring("expected minimum"))

			_, found = routes.Lookup([]string{"HELLO", "WORLD"})
			Expect(found).To(BeTrue())
		})
	})

	When("setting a maximum", func() {
		It("requires that minimum or cannot be found", func() {
			routes := router.MinMaxTokensRouter(0, 1, func(s []string, w io.Writer) error {
				return nil
			})

			writer := &bytes.Buffer{}

			_, found := routes.Lookup([]string{})
			Expect(found).To(BeTrue())

			callback, found := routes.Lookup([]string{"HELLO"})
			Expect(found).To(BeTrue())

			err := callback(nil, writer)
			Expect(err).NotTo(HaveOccurred())
			Expect(writer.String()).NotTo(ContainSubstring("expected maximum"))

			writer.Reset()

			callback, found = routes.Lookup([]string{"HELLO", "WORLD"})
			Expect(found).To(BeFalse())

			err = callback(nil, writer)
			Expect(err).NotTo(HaveOccurred())
			Expect(writer.String()).To(ContainSubstring("expected maximum"))
		})
	})

	When("setting both min and max", func() {
		It("allows the range to be honored", func() {
			routes := router.MinMaxTokensRouter(2, 4, func(s []string, w io.Writer) error {
				return nil
			})

			_, found := routes.Lookup([]string{})
			Expect(found).To(BeFalse())

			_, found = routes.Lookup([]string{"HELLO"})
			Expect(found).To(BeFalse())

			_, found = routes.Lookup([]string{"HELLO", "WORLD"})
			Expect(found).To(BeTrue())

			_, found = routes.Lookup([]string{"HELLO", "WORLD", "HELLO"})
			Expect(found).To(BeTrue())

			_, found = routes.Lookup([]string{"HELLO", "WORLD", "HELLO", "WORLD"})
			Expect(found).To(BeTrue())

			_, found = routes.Lookup([]string{"HELLO", "WORLD", "HELLO", "WORLD", "HELLO"})
			Expect(found).To(BeFalse())
		})
	})
})
