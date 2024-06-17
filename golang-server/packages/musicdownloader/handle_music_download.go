package musicdownloader

import (
	"context"
	"fmt"
	"log"
	"sync"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func handleError(err error, message string) {
	if message == "" {
	  message = "Error making API call"
	}
	if err != nil {
	  log.Printf(message + ": %v", err.Error())
	}
}


// searchYouTube performs a search on YouTube with the given query and filters out playlists.
func searchYouTube(ctx context.Context, query string, maxResults int64, youtubeApiKey string, wg *sync.WaitGroup, resultChan chan <- map[string][]string) {
	defer wg.Done()

	clientOptions := option.WithAPIKey(youtubeApiKey)

	service, err := youtube.NewService(ctx, clientOptions)
	if err != nil {
		resultChan <- map[string][]string{query: nil}
		handleError(err, "Error creating new YouTube client for query")
		return
	}

	// Make the API call to YouTube with the search filter to exclude playlists.
	call := service.Search.List([]string{"id,snippet"}).
		Q(query).
		MaxResults(maxResults).
		Type("video") // This will filter out channels and playlists
	response, err := call.Do()
	if err != nil {
		resultChan <- map[string][]string{query: nil}
		handleError(err, "Error making YouTube API call for query")
        return
	}

	// Collect video URLs
	videoURLs := make([]string, 0)
	for _, item := range response.Items {
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId)
		videoURLs = append(videoURLs, videoURL)
	}

	resultChan <- map[string][]string{query: videoURLs}
}