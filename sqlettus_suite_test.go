package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/antelman107/net-wait-go/wait"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
	"github.com/redis/go-redis/v9"
)

func TestSqlettus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sqlettus Suite")
}

var _ = Describe("CLI", func() {
	It("can start the server", func() {
		port, err := freeport.GetFreePort()
		Expect(err).NotTo(HaveOccurred())

		cli := &CLI{
			Port:     uint(port),
			Filename: "file:test.db?cache=shared&mode=memory",
			Workers:  1,
		}
		go func() {
			defer GinkgoRecover()

			err := cli.Run()
			Expect(err).NotTo(HaveOccurred())
		}()

		ok := wait.New().Do([]string{fmt.Sprintf("localhost:%d", port)})
		Expect(ok).To(BeTrue())

		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%d", port),
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		value, err := client.Ping(context.Background()).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("PONG"))

		value, err = client.Echo(context.Background(), "message").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("message"))

		value, err = client.FlushAll(context.Background()).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("OK"))

		value, err = client.FlushDB(context.Background()).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("OK"))

		value, err = client.Set(context.Background(), "name", "hello", time.Hour).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("OK"))

		value, err = client.Get(context.Background(), "name").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("hello"))

		intVal, err := client.Append(context.Background(), "name", " world").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(intVal).To(BeEquivalentTo(11))

		intVal, err = client.Del(context.Background(), "name").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(intVal).To(BeEquivalentTo(1))

		value, err = client.Get(context.Background(), "name").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal(""))

		intVal, err = client.Decr(context.Background(), "key").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(intVal).To(BeEquivalentTo(-1))

		intVal, err = client.IncrBy(context.Background(), "key", 2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(intVal).To(BeEquivalentTo(1))

		intVal, err = client.DecrBy(context.Background(), "key", 4).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(intVal).To(BeEquivalentTo(-3))

		intVal, err = client.Incr(context.Background(), "key").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(intVal).To(BeEquivalentTo(-2))
	})
})
