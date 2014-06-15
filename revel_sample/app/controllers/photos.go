package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"github.com/yuroyoro/go_shugyo/revel_sample/app/models"
	"github.com/yuroyoro/go_shugyo/revel_sample/app/routes"
)

type Photos struct {
	App
}

func (c Photos) Index(page int) revel.Result {

	records, err := models.LoadPhotos(c.Txn, page)

	if err != nil {
		panic(err)
	}

	fmt.Println(records)

	first := records[0]
	photos := records[1:]

	return c.Render(first, photos)
}

func (c Photos) Save(photo models.Photo) revel.Result {

	photo.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Photos.Index(0))
	}

	err := c.Txn.Insert(&photo)
	if err != nil {
		panic(err)
	}

	return c.Redirect(routes.Photos.Index(0))
}
