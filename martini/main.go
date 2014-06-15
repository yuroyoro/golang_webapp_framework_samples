package main

import (
	"./models"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
	"strconv"
)

type Body struct {
	First  models.Photo
	Photos []models.Photo
}

func loadBody(page int) (*Body, error) {
	photos, err := models.LoadPhotos(page)
	if err != nil {
		return nil, err
	}
	fmt.Println(photos)

	body := Body{
		photos[0],
		photos[1:],
	}

	return &body, nil
}

func main() {
	m := martini.Classic()

	m.Use(render.Renderer())
	m.Use(martini.Static("views"))

	m.Get("/", func(w http.ResponseWriter, r *http.Request, render render.Render) {
		page, err := strconv.Atoi(r.URL.Path[1:])
		if err != nil {
			page = 0
		}

		body, err := loadBody(page)
		if err != nil {
			panic(err)
		}

		render.HTML(200, "index", body)
	})

	m.Post("/", func(w http.ResponseWriter, r *http.Request, render render.Render) {
		url := r.FormValue("url")
		author := r.FormValue("author")

		if url == "" {
			render.Error(500)
			return
		}

		if author == "" {
			render.Error(500)
			return
		}

		fmt.Printf("Save Photo(%s, %s)", url, author)
		photo := models.Photo{
			URL:    url,
			Author: author,
		}
		photo.Save()

		render.Redirect("/", 302)
	})

	m.Run()
}
