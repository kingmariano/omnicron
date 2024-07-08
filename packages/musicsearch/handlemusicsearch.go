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
package musicsearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8000/api/v1/search-song"

// MusicSearchRequest defines the structure of the request sent to the FastAPI server.
type MusicSearchRequest struct {
	Song  string `json:"song"`
	Limit int    `json:"limit,omitempty"`
	Proxy string `json:"proxy,omitempty"`
}

// MusicSearchResponse defines the structure of the JSON response from the FastAPI server.
type MusicSearchResponse struct {
	Tracks struct {
		Hits []struct {
			Type    string `json:"type"`
			Key     string `json:"key"`
			Heading struct {
				Title    string `json:"title"`
				Subtitle string `json:"subtitle"`
			} `json:"heading"`
			Images struct {
				Default string `json:"default"`
				Blurred string `json:"blurred"`
				Play    string `json:"play"`
			} `json:"images"`
			Stores struct {
				Apple struct {
					Actions []struct {
						Type string `json:"type"`
						URI  string `json:"uri"`
					} `json:"actions"`
					Explicit    bool   `json:"explicit"`
					Previewurl  string `json:"previewurl"`
					Coverarturl string `json:"coverarturl"`
					Trackid     string `json:"trackid"`
					Productid   string `json:"productid"`
				} `json:"apple"`
			} `json:"stores"`
			Streams struct {
			} `json:"streams"`
			Artists []struct {
				Alias  string `json:"alias"`
				ID     string `json:"id"`
				Adamid string `json:"adamid"`
			} `json:"artists"`
			Share struct {
				Subject  string `json:"subject"`
				Text     string `json:"text"`
				Href     string `json:"href"`
				Image    string `json:"image"`
				Twitter  string `json:"twitter"`
				HTML     string `json:"html"`
				Avatar   string `json:"avatar"`
				Snapchat string `json:"snapchat"`
			} `json:"share"`
			Alias   string `json:"alias"`
			URL     string `json:"url"`
			Actions []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"actions"`
			Urlparams struct {
				Tracktitle  string `json:"{tracktitle}"`
				Trackartist string `json:"{trackartist}"`
			} `json:"urlparams"`
		} `json:"hits"`
		Next string `json:"next"`
	} `json:"tracks"`
}

// ErrorResponse defines the structure for error JSON responses.
type ErrorResponse struct {
	Detail string `json:"detail"`
}

// FilteredResponse is a struct to store the filtered results it contains the SonName,ShazamURL,SongImage gotten from the shazam API
type FilteredResponse struct {
	SongName string `json:"song_name"`
	ShazamURL string `json:"shazam_url"`
	SongImage string `json:"song_image"`
}

// CallMusicSearchFastAPI makes a request to the FastAPI server endpoint for music search.
func CallMusicSearchFastAPI(request MusicSearchRequest, apiKey string) ([]FilteredResponse, error) {
	// Marshal request data to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		var errorMessage ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorMessage); err != nil {
			return nil, err
		}
		return nil, errors.New(errorMessage.Detail)
	}

	// Decode JSON response
	var response MusicSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	res := response
	// Filter the results based on the share subject criteria
	FiltRes := make([]FilteredResponse, 0, len(res.Tracks.Hits))
    for _, song := range res.Tracks.Hits {
        FiltRes = append(FiltRes, FilteredResponse{
            SongName: song.Share.Subject,
            ShazamURL: song.Share.Href,
            SongImage: song.Share.Image,
        })
    }

   	return FiltRes, nil
}
