package models

import (
	"errors"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel"
)

var DatabaseFile = "photos.db"

type Photo struct {
	Id     int64
	URL    string
	Author string
}

func LoadPhotos(dbm *gorp.Transaction, page int) ([]Photo, error) {
	if page < 0 {
		return nil, errors.New("invalid page number")
	}

	limit := 8
	offset := page * limit

	var photos []Photo

	_, err := dbm.Select(&photos, "SELECT id, url, author FROM photos ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (photo Photo) Validate(v *revel.Validation) {
	v.Required(photo.URL)
	v.Required(photo.Author)
}
