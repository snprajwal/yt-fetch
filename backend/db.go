package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

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

func initDb() error {
	dbUrl := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname= %s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
	)
	var err error
	ytFetchDb, err = sql.Open("postgres", dbUrl)
	return err
}

func insertVideo(res *youtube.SearchResult) error {
	publishedAt, err := time.Parse(time.RFC3339, res.Snippet.PublishedAt)
	if err != nil {
		return err
	}
	_, err = ytFetchDb.Exec(
		`INSERT INTO video VALUES ($1, $2, $3, $4, $5, $6)`,
		res.Id.VideoId,
		res.Snippet.Title,
		res.Snippet.ChannelTitle,
		res.Snippet.Description,
		res.Snippet.Thumbnails.Default.Url,
		publishedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
