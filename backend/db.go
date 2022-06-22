package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"google.golang.org/api/youtube/v3"
)

var (
	ytFetchDb *sql.DB

	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
)

const (
	GET_VIDEOS_UNPAGED = "SELECT * FROM video ORDER BY (id, published_at) DESC LIMIT $1"
	GET_VIDEOS_PAGED   = "SELECT * FROM video WHERE id <= $1 ORDER BY (id, published_at) DESC LIMIT $2"

	SEARCH_VIDEOS_UNPAGED = `SELECT * FROM video WHERE (LOWER(title) LIKE '%' || $1 || '%'
	OR LOWER(description) LIKE '%' || $2 || '%') ORDER BY (id, published_at) DESC LIMIT $3`
	SEARCH_VIDEOS_PAGED = `SELECT * FROM video WHERE (LOWER(title) LIKE '%' || $1 || '%'
	OR LOWER(description) LIKE '%' || $2 || '%') AND id <= $3 ORDER BY (id, published_at) DESC LIMIT $4`
)

func initDb() error {
	dbUrl := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
	)
	var err error
	ytFetchDb, err = sql.Open("postgres", dbUrl)
	return err
}

func bulkInsertVideos(res []*youtube.SearchResult) error {
	txn, err := ytFetchDb.Begin()
	if err != nil {
		return err
	}
	stmt, err := txn.Prepare(pq.CopyIn(
		"video", "slug", "title", "channel", "description", "thumbnail", "published_at",
	))
	if err != nil {
		return err
	}

	for _, r := range res {
		publishedAt, err := time.Parse(time.RFC3339, r.Snippet.PublishedAt)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(
			r.Id.VideoId,
			r.Snippet.Title,
			r.Snippet.ChannelTitle,
			r.Snippet.Description,
			r.Snippet.Thumbnails.Default.Url,
			publishedAt,
		)
		if err != nil {
			return err
		}
	}
	if _, err = stmt.Exec(); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}
	if err = txn.Commit(); err != nil {
		return err
	}
	return nil
}

func getVideos(pageKey int) ([]*ytVideo, error) {
	var rows *sql.Rows
	var err error
	if pageKey == 0 {
		rows, err = ytFetchDb.Query(GET_VIDEOS_UNPAGED, PAGE_SIZE)
	} else {
		rows, err = ytFetchDb.Query(GET_VIDEOS_PAGED, pageKey, PAGE_SIZE)
	}
	if err != nil {
		return nil, err
	}
	videos := []*ytVideo{}

	defer rows.Close()
	for rows.Next() {
		video := ytVideo{}
		err = rows.Scan(&video.id, &video.Slug, &video.Title, &video.Channel,
			&video.Description, &video.Thumbnail, &video.PublishedAt)
		if err != nil {
			return nil, err
		}
		videos = append(videos, &video)
	}

	return videos, nil
}

func searchVideos(query string, pageKey int) ([]*ytVideo, error) {
	var rows *sql.Rows
	var err error
	if pageKey == 0 {
		rows, err = ytFetchDb.Query(SEARCH_VIDEOS_UNPAGED, query, query, PAGE_SIZE)
	} else {
		rows, err = ytFetchDb.Query(SEARCH_VIDEOS_PAGED, query, query, pageKey, PAGE_SIZE)
	}
	if err != nil {
		return nil, err
	}
	videos := []*ytVideo{}

	defer rows.Close()
	for rows.Next() {
		video := ytVideo{}
		err = rows.Scan(&video.id, &video.Slug, &video.Title, &video.Channel,
			&video.Description, &video.Thumbnail, &video.PublishedAt)
		if err != nil {
			return nil, err
		}
		videos = append(videos, &video)
	}

	return videos, nil
}
