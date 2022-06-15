package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const PAGE_SIZE = 10

type ytVideo struct {
	id          int
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Channel     string    `json:"channel"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	PublishedAt time.Time `json:"publishedAt"`
}

type videosRes struct {
	Items      []*ytVideo `json:"items"`
	NextPageId string     `json:"nextPageId,omitempty"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("pong")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	pageKey, ok := r.Context().Value("pageKey").(int)
	if !ok {
		http.Error(w, "[REQ] Failed to read pageId", http.StatusBadRequest)
	}
	res := videosRes{}

	videos, err := getVideos(pageKey)
	if err != nil {
		http.Error(w, "[DB] Failed to get videos", http.StatusInternalServerError)
		return
	}
	if len(videos) == PAGE_SIZE {
		res.NextPageId = strconv.Itoa(videos[PAGE_SIZE-1].id - 1)
	}
	res.Items = videos

	json.NewEncoder(w).Encode(res)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	pageKey, ok := r.Context().Value("pageKey").(int)
	if !ok {
		http.Error(w, "[REQ] Failed to read pageId", http.StatusBadRequest)
	}
	query := r.URL.Query().Get("query")
	res := videosRes{}

	videos, err := searchVideos(query, pageKey)
	if err != nil {
		http.Error(w, "[DB] Failed to get videos", http.StatusInternalServerError)
		return
	}
	if len(videos) == PAGE_SIZE {
		res.NextPageId = strconv.Itoa(videos[PAGE_SIZE-1].id - 1)
	}
	res.Items = videos

	json.NewEncoder(w).Encode(res)
}

func pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageId := r.URL.Query().Get("pageId")
		log.Println("[ID] ", pageId)
		// Use 0 as default value to indicate first page
		pageKey := 0
		if pageId != "" {
			pageKey, _ = strconv.Atoi(pageId)
		}
		log.Println("[KEY] ", pageKey)
		ctx := context.WithValue(r.Context(), "pageKey", pageKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
