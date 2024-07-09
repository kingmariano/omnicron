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
package shazam

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// ShazamResponse defines the structure of the JSON response from the FastAPI server.
type ShazamResponse struct {
	Matches []struct {
		ID            string  `json:"id"`
		Offset        float64 `json:"offset"`
		Timeskew      float64 `json:"timeskew"`
		Frequencyskew float64 `json:"frequencyskew"`
	} `json:"matches"`
	Location struct {
		Accuracy float64 `json:"accuracy"`
	} `json:"location"`
	Timestamp int    `json:"timestamp"`
	Timezone  string `json:"timezone"`
	Track     struct {
		Layout   string `json:"layout"`
		Type     string `json:"type"`
		Key      string `json:"key"`
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Images   struct {
			Background string `json:"background"`
			Coverart   string `json:"coverart"`
			Coverarthq string `json:"coverarthq"`
			Joecolor   string `json:"joecolor"`
		} `json:"images"`
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
		Hub struct {
			Type    string `json:"type"`
			Image   string `json:"image"`
			Actions []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				ID   string `json:"id,omitempty"`
				URI  string `json:"uri,omitempty"`
			} `json:"actions"`
			Options []struct {
				Caption string `json:"caption"`
				Actions []struct {
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"actions"`
				Beacondata struct {
					Type         string `json:"type"`
					Providername string `json:"providername"`
				} `json:"beacondata"`
				Image               string `json:"image"`
				Type                string `json:"type"`
				Listcaption         string `json:"listcaption"`
				Overflowimage       string `json:"overflowimage"`
				Colouroverflowimage bool   `json:"colouroverflowimage"`
				Providername        string `json:"providername"`
			} `json:"options"`
			Explicit    bool   `json:"explicit"`
			Displayname string `json:"displayname"`
		} `json:"hub"`
		Sections []struct {
			Type      string `json:"type"`
			Metapages []struct {
				Image   string `json:"image"`
				Caption string `json:"caption"`
			} `json:"metapages,omitempty"`
			Tabname  string `json:"tabname"`
			Metadata []struct {
				Title string `json:"title"`
				Text  string `json:"text"`
			} `json:"metadata,omitempty"`
			URL string `json:"url,omitempty"`
		} `json:"sections"`
		URL     string `json:"url"`
		Artists []struct {
			Alias  string `json:"alias"`
			ID     string `json:"id"`
			Adamid string `json:"adamid"`
		} `json:"artists"`
		Alias  string `json:"alias"`
		Isrc   string `json:"isrc"`
		Genres struct {
			Primary string `json:"primary"`
		} `json:"genres"`
		Urlparams struct {
			Tracktitle  string `json:"{tracktitle}"`
			Trackartist string `json:"{trackartist}"`
		} `json:"urlparams"`
		Highlightsurls struct {
		} `json:"highlightsurls"`
		Relatedtracksurl string `json:"relatedtracksurl"`
		Albumadamid      string `json:"albumadamid"`
		Trackadamid      string `json:"trackadamid"`
		Releasedate      string `json:"releasedate"`
	} `json:"track"`
	Tagid string `json:"tagid"`
}

type ErrorResponse struct {
	Detail string `json:"detail"`
}

// FilteredResponse is a struct to store the filtered results it contains the SonName,ShazamURL,SongImage gotten from the shazam API
type FilteredResponse struct {
	SongName  string `json:"song_name"`
	ShazamURL string `json:"shazam_url"`
	SongImage string `json:"song_image"`
}

// Calls the "/shazam" endpoint from the fastAPI server
func CallShazamFastAPI(file multipart.File, fileHeader *multipart.FileHeader, apiKey, fastAPIBaseURL string) (*FilteredResponse, error) {
	fastAPIShazamEndpoint := fmt.Sprintf("%s/api/v1/shazam", fastAPIBaseURL) //url to the shazam endpoint in the python server
	// Create a buffer to write our form data to
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Create a form file field
	fw, err := w.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, err
	}
	// Copy the file into the form field
	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}
	// Close the multipart writer to set the terminating boundary
	if err := w.Close(); err != nil {
		return nil, err
	}
	//create the http client
	client := &http.Client{
		Timeout: time.Second * 300,
	}

	req, err := http.NewRequest("POST", fastAPIShazamEndpoint, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Api-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		//handle error when status code is not 200
		var errorMessage ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorMessage); err != nil {
			return nil, err
		}
		return nil, errors.New(errorMessage.Detail)
	}
	var response ShazamResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	res := response
	// Filter the results based on the share subject criteria
	return &FilteredResponse{
		SongName:  res.Track.Share.Subject,
		ShazamURL: res.Track.Share.Href,
		SongImage: res.Track.Share.Image}, nil
}
