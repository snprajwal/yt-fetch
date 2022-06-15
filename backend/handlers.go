package main

import (
	"net/http"
	"time"
)

type ytVideo struct {
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Channel     string    `json:"channel"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	PublishedAt time.Time `json:"publishedAt"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
func listHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
