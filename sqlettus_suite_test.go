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
	var client *redis.Client

	BeforeEach(func() {
		port, err := freeport.GetFreePort()
		Expect(err).NotTo(HaveOccurred())

		cli := &CLI{
			Port:     uint(port),
			Filename: "sqlite://:memory:?cache=shared&mode=memory",
			Workers:  1,
		}
		go func() {
			defer GinkgoRecover()

			err := cli.Run()
			Expect(err).NotTo(HaveOccurred())
		}()

		ok := wait.New().Do([]string{fmt.Sprintf("localhost:%d", port)})
		Expect(ok).To(BeTrue())

		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%d", port),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	})

	It("can send PING", func() {
		value, err := client.Ping(context.Background()).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("PONG"))
	})

	It("can send ECHO", func() {
		value, err := client.Echo(context.Background(), "message").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("message"))
	})

	It("can send FLUSHALL", func() {
		set(client, "hello", "world")
		get(client, "hello", "world")

		value, err := client.FlushAll(context.Background()).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("OK"))

		get(client, "hello", "")
	})

	It("can send FLUSHDB", func() {
		set(client, "hello", "world")
		get(client, "hello", "world")

		value, err := client.FlushDB(context.Background()).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("OK"))

		get(client, "hello", "")
	})

	It("can send SET", func() {
		value, err := client.Set(context.Background(), "mykey", "Hello", time.Hour).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("OK"))

		get(client, "mykey", "Hello")
	})

	It("can send GET", func() {
		set(client, "mykey", "Hello")

		value, err := client.Get(context.Background(), "mykey").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("Hello"))
	})

	It("can send APPEND", func() {
		// Add EXISTS check
		value, err := client.Append(context.Background(), "mykey", "Hello").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(5))

		value, err = client.Append(context.Background(), "mykey", " World").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(11))

		get(client, "mykey", "Hello World")
	})

	It("can send DEL", func() {
		set(client, "key1", "Hello")
		set(client, "key2", "World")

		value, err := client.Del(context.Background(), "key1", "key2", "key3").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(2))

		get(client, "key1", "")
		get(client, "key2", "")
	})

	It("can send DECR", func() {
		set(client, "mykey", "10")

		value, err := client.Decr(context.Background(), "mykey").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(9))

		set(client, "mykey", "234293482390480948029348230948")

		value, err = client.Decr(context.Background(), "mykey").Result()
		Expect(err).To(HaveOccurred())
		Expect(value).To(BeEquivalentTo(0))

		value, err = client.Decr(context.Background(), "newkey").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(-1))

		get(client, "newkey", "-1")
	})

	It("can send DECRBY", func() {
		set(client, "mykey", "10")

		value, err := client.DecrBy(context.Background(), "mykey", 3).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(7))

		get(client, "mykey", "7")

		set(client, "mykey", "234293482390480948029348230948")

		value, err = client.DecrBy(context.Background(), "mykey", 3).Result()
		Expect(err).To(HaveOccurred())
		Expect(value).To(BeEquivalentTo(0))

		value, err = client.DecrBy(context.Background(), "newkey", 3).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(-3))

		get(client, "newkey", "-3")
	})

	It("can send INCR", func() {
		set(client, "mykey", "10")

		value, err := client.Incr(context.Background(), "mykey").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(11))

		get(client, "mykey", "11")

		set(client, "mykey", "234293482390480948029348230948")

		value, err = client.Incr(context.Background(), "mykey").Result()
		Expect(err).To(HaveOccurred())
		Expect(value).To(BeEquivalentTo(0))

		value, err = client.Incr(context.Background(), "newkey").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(1))

		get(client, "newkey", "1")
	})

	It("can send INCRBY", func() {
		set(client, "mykey", "10")

		value, err := client.IncrBy(context.Background(), "mykey", 5).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(15))

		get(client, "mykey", "15")

		set(client, "mykey", "234293482390480948029348230948")

		value, err = client.IncrBy(context.Background(), "mykey", 1).Result()
		Expect(err).To(HaveOccurred())
		Expect(value).To(BeEquivalentTo(0))

		value, err = client.IncrBy(context.Background(), "newkey", 3).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(3))

		get(client, "newkey", "3")
	})

	It("can send INCRBYFLOAT", func() {
		set(client, "mykey", "10.50")

		value, err := client.IncrByFloat(context.Background(), "mykey", 0.1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(10.6))

		value, err = client.IncrByFloat(context.Background(), "mykey", -5).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(5.6))

		set(client, "mykey", "5.0e3")

		value, err = client.IncrByFloat(context.Background(), "mykey", 2.0e2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(5200))
	})

	It("can send GETDEL", func() {
		set(client, "mykey", "Hello")

		value, err := client.GetDel(context.Background(), "mykey").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("Hello"))

		get(client, "mykey", "")
	})

	It("can send GETRANGE", func() {
		set(client, "mykey", "This is a string")

		value, err := client.GetRange(context.Background(), "mykey", 0, 3).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("This"))

		value, err = client.GetRange(context.Background(), "mykey", -3, -1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("ing"))

		value, err = client.GetRange(context.Background(), "mykey", 0, -1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("This is a string"))

		value, err = client.GetRange(context.Background(), "mykey", 10, 100).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("string"))
	})

	It("can send MSET", func() {
		value, err := client.MSet(context.TODO(), "key1", "value1", "key2", "value2").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("OK"))

		get(client, "key1", "value1")
		get(client, "key2", "value2")
	})

	It("can send MGET", func() {
		set(client, "key1", "Hello")
		set(client, "key2", "World")

		values, err := client.MGet(context.TODO(), "key1", "key2", "nonexisting").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(Equal([]interface{}{
			"Hello",
			"World",
			nil,
		}))
	})

	It("can send UNLINK", func() {
		set(client, "key1", "Hello")
		set(client, "key2", "World")

		value, err := client.Unlink(context.TODO(), "key1", "key2", "key3").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(2))

		get(client, "key1", "")
		get(client, "key2", "")
	})

	It("can send STRLEN", func() {
		set(client, "mykey", "Hello world")

		value, err := client.StrLen(context.TODO(), "mykey").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(11))

		value, err = client.StrLen(context.TODO(), "nonexisting").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(0))
	})

	It("can send RPUSHX", func() {
		value, err := client.RPush(context.TODO(), "mylist", "Hello").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(1))

		value, err = client.RPushX(context.TODO(), "mylist", "World").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(2))

		value, err = client.RPushX(context.TODO(), "myotherlist", "World").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(0))

		values, err := client.LRange(context.TODO(), "mylist", 0, -1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(Equal([]string{"Hello", "World"}))

		values, err = client.LRange(context.TODO(), "myotherlist", 0, -1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(BeEmpty())
	})

	It("can send RPUSH", func() {
		value, err := client.RPush(context.TODO(), "mylist", "hello").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(1))

		value, err = client.RPushX(context.TODO(), "mylist", "world").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(2))

		values, err := client.LRange(context.TODO(), "mylist", 0, -1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(Equal([]string{"hello", "world"}))
	})

	It("can send LRANGE", func() {
		value, err := client.RPush(context.TODO(), "mylist", "one").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(1))

		value, err = client.RPush(context.TODO(), "mylist", "two").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(2))

		value, err = client.RPush(context.TODO(), "mylist", "three").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo(3))

		values, err := client.LRange(context.TODO(), "mylist", 0, 0).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(Equal([]string{"one"}))

		values, err = client.LRange(context.TODO(), "mylist", -3, 2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(Equal([]string{"one", "two", "three"}))

		values, err = client.LRange(context.TODO(), "mylist", -100, 100).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(Equal([]string{"one", "two", "three"}))

		values, err = client.LRange(context.TODO(), "mylist", 5, 10).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(BeEmpty())

		values, err = client.LRange(context.TODO(), "nonexisting", 0, 1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(BeEmpty())
	})

	// It("can start the server", func() {

	// 	intVal, err = client.RPushX(context.TODO(), "mykey", "two").Result()
	// 	Expect(err).To(HaveOccurred())
	// 	Expect(intVal).To(BeEquivalentTo(0))
	// })
})

func set(client *redis.Client, key, value string) {
	err := client.Set(context.Background(), key, value, time.Hour).Err()
	Expect(err).NotTo(HaveOccurred())
}

func get(client *redis.Client, key, expected string) {
	actual, err := client.Get(context.Background(), key).Result()
	if expected == "" {
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("redis: nil"))
	} else {
		Expect(err).NotTo(HaveOccurred())
		Expect(expected).To(Equal(actual))
	}
}
