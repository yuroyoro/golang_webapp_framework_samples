package controllers

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	r "github.com/revel/revel"
	"github.com/revel/revel/modules/db/app"
	m "github.com/yuroyoro/go_shugyo/revel_sample/app/models"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}

	Dbm.AddTableWithName(m.Photo{}, "photos").SetKeys(true, "Id")

	Dbm.TraceOn("[gorp]", r.INFO)

	err := Dbm.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}

	photos := []*m.Photo{
		&m.Photo{URL: "http://25.media.tumblr.com/88b812f5f9c3d7b83560fd635435d538/tumblr_mx3tlblmY21st5lhmo1_1280.jpg", Author: "Yinan Chen"},
		&m.Photo{URL: "http://25.media.tumblr.com/95c842c76d60b7bc982d92c76216d037/tumblr_mx3tnm96k81st5lhmo1_1280.jpg", Author: "Thanun Buranapong"},
		&m.Photo{URL: "http://24.media.tumblr.com/c35afcc83e18ea7875160f64c039f471/tumblr_mwhdohfePJ1st5lhmo1_1280.jpg", Author: "Linh Nguyen"},
		&m.Photo{URL: "http://24.media.tumblr.com/e100564a3e73c9456acddb9f62f96c79/tumblr_mufs8mix841st5lhmo1_1280.jpg", Author: "Charlie Foster"},
		&m.Photo{URL: "http://24.media.tumblr.com/ac840897b5f73fa6bc43f73996f02572/tumblr_mrraat0H431st5lhmo1_1280.jpg", Author: "Rula Sibai"},
		&m.Photo{URL: "http://24.media.tumblr.com/23e3f4bb271b8bdc415275fb7061f204/tumblr_mve3rvxwaP1st5lhmo1_1280.jpg", Author: "Jonas Nilsson Lee"},
		&m.Photo{URL: "http://31.media.tumblr.com/aa1779a718c2844969f23c4f5dec86b1/tumblr_mvyxhonf601st5lhmo1_1280.jpg", Author: "Linh Nguyen"},
		&m.Photo{URL: "http://24.media.tumblr.com/d6b9403c704c3e5aa1725c106e8a9430/tumblr_mvyxd9PUpZ1st5lhmo1_1280.jpg", Author: "Dillon McIntosh"},
	}

	for _, photo := range photos {
		if err := Dbm.Insert(photo); err != nil {
			panic(err)
		}
	}

}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
