package models

import (
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
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

	if len(photos) != 1 {
		t.Error("Photo.Save() failed")
	}

	if photos[0].URL != photo.URL || photos[0].Author != photo.Author {
		t.Error("Photo.Save() failed")
	}

}

func TestLoadPhotos(t *testing.T) {
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

	// when given page 0
	photos, err := LoadPhotos(0)

	if len(photos) != 8 {
		t.Errorf("Photo.LoadPhotos(0) expected to return 8 records, but %d", len(photos))
	}

	first := photos[0]
	if first.URL != "http://example.com/19" || first.Author != "author_19" {
		fmt.Println(first)
		t.Error("Photo.LoadPhotos(0) returns unexpected reocreds")
	}

	// when given page 2
	photos, err = LoadPhotos(2)

	if len(photos) != 4 {
		t.Errorf("Photo.LoadPhotos(2) expected to return 4 records, but %d", len(photos))
	}

	last := photos[len(photos)-1]
	if last.URL != "http://example.com/0" || last.Author != "author_0" {
		fmt.Println(last)
		t.Error("Photo.LoadPhotos(2) returns unexpected reocreds")
	}

	// when given page 99
	photos, err = LoadPhotos(99)
	if len(photos) != 0 {
		t.Errorf("Photo.LoadPhotos(99) expected to return empty records, but %d", len(photos))
	}

}
