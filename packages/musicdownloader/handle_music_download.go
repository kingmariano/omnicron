// Copyright (c) 2024 Charles Ozochukwu

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package musicdownloader

import (
	"context"
	"fmt"
	"github.com/kingmariano/omnicron/packages/videodownloader"
	"github.com/kingmariano/omnicron/utils"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// searchYouTube performs a search for audio music on YouTube with the given query and filters out playlists.
func SearchMusicOnYouTube(ctx context.Context, query string, maxResults int64, youtubeAPIKey, _ string) ([]string, error) {
	clientOptions := option.WithAPIKey(youtubeAPIKey)
	service, err := youtube.NewService(ctx, clientOptions)
	if err != nil {
		return []string{}, fmt.Errorf("error creating new YouTube client for query %v: %v", query, err)
	}

	// Make the API call to YouTube with the search filter to exclude playlists.
	call := service.Search.List([]string{"id,snippet"}).
		Q(query).
		MaxResults(maxResults).
		Type("video") // This will filter out channels and playlists
	response, err := call.Do()
	if err != nil {
		return []string{}, fmt.Errorf("failed to make request to youtube client %v", err)
	}

	// Collect video URLs
	videoURLs := make([]string, 0)
	for _, item := range response.Items {
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId)
		videoURLs = append(videoURLs, videoURL)
	}
	return videoURLs, nil
}

func downloadYoutubeLinkAndConvertToMp3(ctx context.Context, query string, maxResults int64, youtubeAPIKey, cloudinaryURL, outputPath string) ([]string, error) {
	// Perform a search on YouTube for the given query and retrieve video URLs
	//set output path to where the song will be downloaded

	urlList, err := SearchMusicOnYouTube(ctx, query, maxResults, youtubeAPIKey, cloudinaryURL)
	if err != nil {
		return []string{}, err
	}

	// If no video URLs are found, return an empty list
	if len(urlList) == 0 {
		return []string{}, nil
	}

	// Download all the video in the list
	videoPathList := make([]string, 0)
	for _, url := range urlList {
		videopath, err := videodownloader.DownloadVideoData(url, utils.OutputName, outputPath, "")
		if err != nil {
			return []string{}, err
		}
		videoPathList = append(videoPathList, videopath)
	}

	// Convert the downloaded videos to MP3 format
	audioPathList := make([]string, 0)
	for _, videopath := range videoPathList {
		audiopath, err := utils.ConvertFileToMP3(videopath)
		if err != nil {
			return []string{}, err
		}
		audioPathList = append(audioPathList, audiopath)
	}

	// Upload the converted audio file to Cloudinary and retrieve direct URLs
	audioFileURLList := make([]string, 0)
	for _, audiopath := range audioPathList {
		audioDirectURL, err := utils.HandleFileUpload(ctx, audiopath, cloudinaryURL)
		if err != nil {
			return []string{}, err
		}
		audioFileURLList = append(audioFileURLList, audioDirectURL)
	}

	// Return the list of direct URLs to the uploaded audio files on Cloudinary
	return audioFileURLList, nil
}
