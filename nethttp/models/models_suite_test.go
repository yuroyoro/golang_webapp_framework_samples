package models_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/coopernurse/gorp"
	"github.com/yuroyoro/go_shugyo/nethttp/models"
	"testing"
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = Describe("Photo", func() {
	var (
		dbmap *gorp.DbMap
		err   error
	)

	BeforeEach(func() {
		models.DatabaseFile = "photos_test.db"
		dbmap, err = models.InitDb()
		if err != nil {
			panic(err)
		}

		err = dbmap.TruncateTables()
		if err != nil {
			panic(err)
		}
	})

	AfterEach(func() {
		dbmap.Db.Close()
	})

	Describe("Insert", func() {
		It("should inserts new record", func() {
			photo := models.Photo{
				URL:    "http://example.com",
				Author: "yuroyoro",
			}

			photo.Save()

			var photos []models.Photo

			_, err = dbmap.Select(&photos, "SELECT id, url, author FROM photos ORDER BY id ASC ")
			if err != nil {
				Fail("Failed to load records from database")
			}

			Expect(len(photos)).To(Equal(1))
			Expect(photos[0].URL).To(Equal(photo.URL))
			Expect(photos[0].Author).To(Equal(photo.Author))

		})
	})

	Describe("LoadPhotos", func() {
		BeforeEach(func() {
			for i := 0; i < 20; i++ {
				photo := models.Photo{
					URL:    fmt.Sprintf("http://example.com/%d", i),
					Author: fmt.Sprintf("author_%d", i),
				}

				photo.Save()
			}
		})

		Context("when given page 0", func() {
			It("should returns first page", func() {
				photos, err := models.LoadPhotos(0)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(photos)).To(Equal(8))

				first := photos[0]

				Expect(first.URL).To(Equal("http://example.com/19"))
				Expect(first.Author).To(Equal("author_19"))
			})
		})

		Context("when given page 2", func() {
			It("should returns last page", func() {
				photos, err := models.LoadPhotos(2)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(photos)).To(Equal(4))

				last := photos[len(photos)-1]

				Expect(last.URL).To(Equal("http://example.com/0"))
				Expect(last.Author).To(Equal("author_0"))
			})
		})

		Context("when given page 99", func() {
			It("should returns empty", func() {
				photos, err := models.LoadPhotos(99)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(photos)).To(Equal(0))
			})
		})

		Context("when given page -1", func() {
			It("should returns error", func() {
				photos, err := models.LoadPhotos(-1)

				Expect(err).Should(HaveOccurred())
				Expect(photos).To(BeNil())
			})
		})

	})
})
