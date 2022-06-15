package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/speps/go-hashids"
)

const PAGE_SIZE = 10

var h *hashids.HashID

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

func initHasher() error {
	var err error
	h, err = hashids.New()
	return err
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
		nextPageId, err := h.Encode([]int{videos[PAGE_SIZE-1].id - 1})
		if err != nil {
			http.Error(w, "[HASH] Failed to encode pageId", http.StatusInternalServerError)
			return
		}
		res.NextPageId = nextPageId
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
		nextPageId, err := h.Encode([]int{videos[PAGE_SIZE-1].id - 1})
		if err != nil {
			http.Error(w, "[HASH] Failed to encode pageId", http.StatusInternalServerError)
			return
		}
		res.NextPageId = nextPageId
	}
	res.Items = videos

	json.NewEncoder(w).Encode(res)
}

func pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageId := r.URL.Query().Get("pageId")
		// Use 0 as default value to indicate first page
		pageKey := 0
		if pageId != "" {
			hash, err := h.DecodeWithError(pageId)
			if err != nil {
				http.Error(w, "[HASH] Failed to decode pageId", http.StatusBadRequest)
				return
			}
			pageKey = hash[0]
		}
		ctx := context.WithValue(r.Context(), "pageKey", pageKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
