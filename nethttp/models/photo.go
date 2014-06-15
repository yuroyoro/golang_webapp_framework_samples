package models

import (
	"database/sql"
	"errors"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

var DatabaseFile = "photos.db"

type Photo struct {
	Id     int64
	URL    string
	Author string
}

func (p Photo) Save() error {
	dbmap, err := InitDb()
	if err != nil {
		return err
	}

	defer dbmap.Db.Close()

	dbmap.Insert(&p)
	if err != nil {
		return err
	}

	return nil
}

func LoadPhotos(page int) ([]Photo, error) {
	dbmap, err := InitDb()
	if err != nil {
		return nil, err
	}

	defer dbmap.Db.Close()

	if page < 0 {
		return nil, errors.New("invalid page number")
	}

	limit := 8
	offset := page * limit

	var photos []Photo

	_, err = dbmap.Select(&photos, "SELECT id, url, author FROM photos ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	return photos, nil

}

func InitDb() (*gorp.DbMap, error) {
	db, err := sql.Open("sqlite3", DatabaseFile)
	if err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Photo{}, "photos").SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	return dbmap, nil
}
