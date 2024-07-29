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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kingmariano/omnicron/packages/videodownloader"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
	"time"
)

func CallSearchYoutubeFastdownloadYoutubeLink(ctx context.Context, request SongRequest, apiKey, fastAPIBaseURL, outputPath, cloudinaryURL string) (string, error) {
	fastAPISearchYoutubeEndPoint := fmt.Sprintf("%s/api/v1/search_youtube", fastAPIBaseURL) //url to the search_youtube endpoint in the FastAPI server
	// marshal the request to json format for sending to the server
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Timeout: time.Second * 300,
	}
	req, err := http.NewRequest("POST", fastAPISearchYoutubeEndPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//handle error when status code is not 200
	if resp.StatusCode != 200 {
		var errorMessage ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorMessage); err != nil {
			return "", err
		}
		return "", errors.New(errorMessage.Detail)
	}
	var response SongResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	if response.Response == "" {
		return "", errors.New("no results found")
	}
	// Download all the video in the list
	videopath, err := videodownloader.DownloadVideoData(response.Response, utils.OutputName, outputPath, "")
	if err != nil {
		return "", err
	}
	// Convert the downloaded videos to MP3 format
	audiopath, err := utils.ConvertFileToMP3(videopath)
	if err != nil {
		return "", err
	}
	// Upload the converted audio file to Cloudinary and retrieve direct URLs
	audioDirectURL, err := utils.HandleFileUpload(ctx, audiopath, cloudinaryURL)
	if err != nil {
		return "", err
	}
	// Return the direct URL to the uploaded audio file on Cloudinary
	return audioDirectURL, nil
}
