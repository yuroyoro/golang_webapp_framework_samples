package models

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConverySave(t *testing.T) {
	Convey("Insert", t, func() {
		DatabaseFile = "photos_test.db"
		dbmap, err := InitDb()

		if err != nil {
			t.FailNow()
		}

		defer dbmap.Db.Close()

		err = dbmap.TruncateTables()

		if err != nil {
			t.FailNow()
		}

		photo := Photo{
			URL:    "http://example.com",
			Author: "yuroyoro",
		}

		photo.Save()

		var photos []Photo

		_, err = dbmap.Select(&photos, "SELECT id, url, author FROM photos ORDER BY id ASC ")
		if err != nil {
			t.FailNow()
		}

		So(len(photos), ShouldEqual, 1)
		So(photos[0].URL, ShouldEqual, photo.URL)
		So(photos[0].Author, ShouldEqual, photo.Author)
	})

}

func TestConveryLoadPhotos(t *testing.T) {
	Convey("LoadPhotos", t, func() {

		DatabaseFile = "photos_test.db"
		dbmap, err := InitDb()

		if err != nil {
			t.FailNow()
		}

		defer dbmap.Db.Close()

		err = dbmap.TruncateTables()

		// insert test data
		for i := 0; i < 20; i++ {
			photo := Photo{
				URL:    fmt.Sprintf("http://example.com/%d", i),
				Author: fmt.Sprintf("author_%d", i),
			}

			photo.Save()
		}

		Convey("when given page 0", func() {
			photos, _ := LoadPhotos(0)

			So(len(photos), ShouldEqual, 8)

			first := photos[0]

			So(first.URL, ShouldEqual, "http://example.com/19")
			So(first.Author, ShouldEqual, "author_19")

		})

		Convey("when given page 2", func() {
			photos, _ := LoadPhotos(2)

			So(len(photos), ShouldEqual, 4)

			last := photos[len(photos)-1]

			So(last.URL, ShouldEqual, "http://example.com/0")
			So(last.Author, ShouldEqual, "author_0")

		})

		Convey("when given page 99", func() {
			photos, _ := LoadPhotos(99)

			So(len(photos), ShouldEqual, 0)
		})

		Convey("when given page -1", func() {
			photos, err := LoadPhotos(-1)

			So(err.Error(), ShouldEqual, "invalid page number")
			So(photos, ShouldBeNil)
		})
	})
}
