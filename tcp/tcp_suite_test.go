package tcp_test

import (
	"testing"

	"github.com/jtarchie/sqlettus/tcp"
	"github.com/jtarchie/sqlettus/tcp/handlers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

func TestTcp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TCP Suite")
}

func startServer(handler tcp.Handler) (int, *tcp.Server) {
	port, err := freeport.GetFreePort()
	Expect(err).NotTo(HaveOccurred())

	server, err := tcp.NewServer(port, 1)
	Expect(err).NotTo(HaveOccurred())

	go func() {
		defer GinkgoRecover()

		err := server.Listen(handler)
		Expect(err).NotTo(HaveOccurred())
	}()

	return port, server
}

var _ = Describe("When starting a TCP Server", func() {
	It("accepts a connection", func() {
		port, server := startServer(&handlers.Echo{})
		defer server.Close()

		response, err := tcp.Write(port, "echo\r\n")
		Expect(err).NotTo(HaveOccurred())
		Expect(response).To(Equal("echo"))
	})

	When("the handler errors on the client", func() {
		It("server continues accepting connections", func() {
			_, server := startServer(&handlers.Error{})
			defer server.Close()

			// _, err := tcp.Write(port, "echo\r\n")
			// Expect(err).To(HaveOccurred())
		})
	})
})
