package db_test

import (
	"context"

	"github.com/jtarchie/sqlettus/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	var client *db.Client

	BeforeEach(func() {
		var err error

		client, err = db.NewClient("sqlite://:memory:?cache=shared&mode=memory")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		client.Close()
	})

	Describe("ListSet", func() {
		It("returns the key was not found", func() {
			found, err := client.ListSet(context.Background(), "key", 0, "123")
			Expect(found).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})

		It("sets a value a position", func() {
			_, _, _ = client.ListRightPushUpsert(context.Background(), "mylist", "one")
			_, _, _ = client.ListRightPushUpsert(context.Background(), "mylist", "two")
			_, _, _ = client.ListRightPushUpsert(context.Background(), "mylist", "three")

			values, _ := client.ListRange(context.Background(), "mylist", 0, -1)
			Expect(values).To(Equal([]string{"one", "two", "three"}))

			found, err := client.ListSet(context.Background(), "mylist", 0, "four")
			Expect(found).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())

			found, err = client.ListSet(context.Background(), "mylist", -2, "five")
			Expect(found).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())

			values, err = client.ListRange(context.Background(), "mylist", 0, -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal([]string{"four", "five", "three"}))
		})
	})

	Describe("ListLength", func() {
		It("returns a count of number of elements", func() {
			length, err := client.ListLength(context.Background(), "mylist")
			Expect(err).NotTo(HaveOccurred())
			Expect(length).To(BeEquivalentTo(0))

			_, _, _ = client.ListRightPushUpsert(context.Background(), "mylist", "Hello")
			_, _, _ = client.ListRightPushUpsert(context.Background(), "mylist", "World")

			length, err = client.ListLength(context.Background(), "mylist")
			Expect(err).NotTo(HaveOccurred())
			Expect(length).To(BeEquivalentTo(2))
		})

		When("value is not an array", func() {
			It("returns an error", func() {
				_ = client.Set(context.Background(), "notlist", "string")
				length, err := client.ListLength(context.Background(), "notlist")
				Expect(err).To(HaveOccurred())
				Expect(length).To(BeEquivalentTo(0))
			})
		})
	})

	Describe("ListInsert", func() {
		It("inserts values at a pivot point", func() {
			_, _, _ = client.ListRightPushUpsert(context.Background(), "mylist", "Hello")
			_, _, _ = client.ListRightPushUpsert(context.Background(), "mylist", "World")

			index, found, err := client.ListInsert(context.Background(), "mylist", -1, "World", "There")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(index).To(BeEquivalentTo(3))

			values, err := client.ListRange(context.Background(), "mylist", 0, -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal([]string{"Hello", "There", "World"}))

			index, found, err = client.ListInsert(context.Background(), "mylist", 1, "World", "Greetings")
			Expect(found).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
			Expect(index).To(BeEquivalentTo(4))

			values, err = client.ListRange(context.Background(), "mylist", 0, -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal([]string{"Hello", "There", "World", "Greetings"}))
		})

		It("returns not found when key does not exist", func() {
			index, found, err := client.ListInsert(context.Background(), "mylist", -1, "a", "b")
			Expect(found).To(BeFalse())
			Expect(err).NotTo(HaveOccurred())
			Expect(index).To(BeEquivalentTo(0))
		})
	})

	Describe("ListRange", func() {
		It("handles zero index and negative indices", func() {
			_, _, _ = client.ListRightPushUpsert(context.Background(), "mylist", "one")
			_, _, _ = client.ListRightPushUpsert(context.Background(), "mylist", "two")
			_, _, _ = client.ListRightPushUpsert(context.Background(), "mylist", "three")

			values, err := client.ListRange(context.Background(), "mylist", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal([]string{"one"}))

			values, err = client.ListRange(context.Background(), "mylist", -3, 2)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal([]string{"one", "two", "three"}))

			values, err = client.ListRange(context.Background(), "mylist", -100, 100)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal([]string{"one", "two", "three"}))

			values, err = client.ListRange(context.Background(), "mylist", 5, 10)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(BeEmpty())
		})

		It("reports missing keys", func() {
			values, err := client.ListRange(context.Background(), "mylist", 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(BeEmpty())
		})
	})

	Describe("ListRightPushUpsert", func() {
		It("returns the index of the value pushed", func() {
			index, found, err := client.ListRightPushUpsert(context.Background(), "mylist", "hello")
			Expect(index).To(BeEquivalentTo(1))
			Expect(found).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())

			index, found, err = client.ListRightPushUpsert(context.Background(), "mylist", "world")
			Expect(index).To(BeEquivalentTo(2))
			Expect(found).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())

			values, err := client.ListRange(context.Background(), "mylist", 0, -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal([]string{"hello", "world"}))
		})

		When("the key already exists of a different type", func() {
			It("returns found, but does not push", func() {
				_ = client.Set(context.Background(), "notlist", "string")

				index, found, err := client.ListRightPushUpsert(context.Background(), "notlist", "hello")
				Expect(index).To(BeEquivalentTo(0))
				Expect(found).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("ListRightPush", func() {
		When("key does not already exist", func() {
			It("returns the index ", func() {
				index, err := client.ListRightPush(context.Background(), "mylist", "hello")
				Expect(index).To(BeEquivalentTo(0))
				Expect(err).ToNot(HaveOccurred())

				values, err := client.ListRange(context.Background(), "mylist", 0, -1)
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(BeEmpty())
			})
		})

		When("key already exist", func() {
			When("value is a different type", func() {
				It("returns found, but does not push", func() {
					_ = client.Set(context.Background(), "notlist", "string")

					index, err := client.ListRightPush(context.Background(), "notlist", "hello")
					Expect(err).To(HaveOccurred())
					Expect(index).To(BeEquivalentTo(0))
				})
			})

			It("returns the index", func() {
				index, found, err := client.ListRightPushUpsert(context.Background(), "mylist", "hello")
				Expect(index).To(BeEquivalentTo(1))
				Expect(found).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())

				index, err = client.ListRightPush(context.Background(), "mylist", "hello")
				Expect(index).To(BeEquivalentTo(2))
				Expect(found).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())

				values, err := client.ListRange(context.Background(), "mylist", 0, -1)
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(Equal([]string{"hello", "hello"}))
			})
		})
	})
})
