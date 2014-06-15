package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

type Photo struct {
	Id     int64
	URL    string
	Author string
}

func (p Photo) Save() error {
	dbmap, err := initDb()
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
	dbmap, err := initDb()
	if err != nil {
		return nil, err
	}

	defer dbmap.Db.Close()

	limit := 8
	offset := page * limit

	var photos []Photo

	_, err = dbmap.Select(&photos, "SELECT url, author FROM photos ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	return photos, nil

}

func initDb() (*gorp.DbMap, error) {
	db, err := sql.Open("sqlite3", "photos.db")
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
