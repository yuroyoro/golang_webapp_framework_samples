package main

import (
	"fmt"
	"github.com/yuroyoro/go_shugyo/nethttp/models"
	"html/template"
	"net/http"
	"strconv"
)

var indexTemplate *template.Template

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Path[1:])
	if err != nil {
		page = 0
	}

	body, err := loadBody(page)
	if err != nil {
		panic(err)
	}

	indexTemplate.Execute(w, body)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("url")
	author := r.FormValue("author")

	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if author == "" {
		http.Error(w, "Twitted ID is required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Save Photo(%s, %s)", url, author)
	photo := models.Photo{
		URL:    url,
		Author: author,
	}
	photo.Save()

	http.Redirect(w, r, "/", http.StatusFound)
}

func init() {
	indexTemplate = template.Must(template.ParseFiles("views/index.html"))
}

func main() {
	http.Handle("/css/layouts/", http.StripPrefix("/css/layouts/", http.FileServer(http.Dir("views/css/layouts"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveHandler)
	http.ListenAndServe(":5050", nil)
}
