package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var YT_API_KEY = os.Getenv("YT_API_KEY")

func fetchLatestVideos(query string) {
	ytSvc, err := youtube.NewService(context.Background(), option.WithAPIKey(YT_API_KEY))
	if err != nil {
		log.Fatal("[YT API ERROR] ", err)
	}

	for {
		// Get videos published in last 30 seconds
		publishedAfterTime := time.Now().
			Add(-30 * time.Second).
			Format(time.RFC3339)

		/* YouTube has an average upload frequency of
		500 hours of video content per minute.
		Assuming an average video length of 5 min,
		this results in a maximum of 100 videos per minute,
		or 50 videos per 30s. By using 100 as the limit for
		maxResults in the query, we ensure that the response
		is not paginated from YouTube's side. */
		req := ytSvc.Search.
			List([]string{"snippet"}).
			MaxResults(100).
			Order("date").
			PublishedAfter(publishedAfterTime).
			Q(query).
			Type("video")

		response, err := req.Do()
		if err != nil {
			log.Println("[FETCH ERROR] ", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, res := range response.Items {
			if err = insertVideo(res); err != nil {
				log.Println("[DB ERROR] ", err)
				continue
			}
		}
		time.Sleep(30 * time.Second)
	}
}
