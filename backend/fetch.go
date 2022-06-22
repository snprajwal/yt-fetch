package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	YT_API_KEY = os.Getenv("YT_API_KEY")
	ytSvc      *youtube.Service
)

func fetchLatestVideos(query string) {
	var err error
	ytSvc, err = youtube.NewService(context.Background(), option.WithAPIKey(YT_API_KEY))
	if err != nil {
		log.Fatal("[YT API ERROR] ", err)
	}
	log.Println("[YT API] Service creation successful")

	// Initially populate DB with videos from last 30 min
	fetchSeed(query)

	for {
		// Get videos published in last 30 seconds
		publishedAfterTime := time.Now().
			Add(-30 * time.Second).
			Format(time.RFC3339)

		req := ytSvc.Search.
			List([]string{"snippet"}).
			MaxResults(50).
			Order("date").
			PublishedAfter(publishedAfterTime).
			Q(query).
			Type("video")

		res, err := req.Do()
		if err != nil {
			log.Println("[YT API ERROR] ", err)
			time.Sleep(30 * time.Second)
			continue
		}
		log.Println("[YT API] Request successful")

		if err = bulkInsertVideos(res.Items); err != nil {
			log.Println("[DB ERROR] ", err)
			continue
		}
		time.Sleep(30 * time.Second)
	}
}

func fetchSeed(query string) {
	publishedAfterTime := time.Now().
		Add(-30 * time.Minute).
		Format(time.RFC3339)

	req := ytSvc.Search.
		List([]string{"snippet"}).
		MaxResults(50).
		Order("date").
		PublishedAfter(publishedAfterTime).
		Q(query).
		Type("video")

	res, err := req.Do()
	if err != nil {
		log.Println("[YT API ERROR] ", err)
		return
	}
	log.Println("[YT API] Request successful")

	if err = bulkInsertVideos(res.Items); err != nil {
		log.Println("[DB ERROR] ", err)
		return
	}
	log.Println("[DB] Insert successful")
}
