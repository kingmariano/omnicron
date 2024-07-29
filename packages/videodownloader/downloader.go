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
package videodownloader

import (
	"encoding/json"
	"fmt"
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

type DownloadParams struct {
	URL        string `json:"url"`
	Resolution string `json:"resolution"`
}

// DownloadVideo handles the video download process.
// It accepts a POST request with JSON body containing URL and resolution.
// It downloads the video from the provided URL, stores it in a temporary directory,
// uploads the video to Cloudinary, and returns the Cloudinary URL of the uploaded video.
//
// Parameters:
//  w http.ResponseWriter: The response writer for the HTTP request.
//  r *http.Request: The HTTP request.
//  cfg *config.ApiConfig: The API configuration.
//
// Return values:
//  None.
type ResponseMsg struct {
	Response string `json:"response"`
}

func DownloadVideo(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {
	ctx := r.Context()
	decode := json.NewDecoder(r.Body)
	params := DownloadParams{}
	err := decode.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling json, %v", err))
		return
	}
	//creates a temporary file to store the downloaded video
	folderPath, err := utils.CreateUniqueFolder(utils.BasePath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	videoPath, err := DownloadVideoData(params.URL, utils.OutputName, folderPath, params.Resolution)
	if err != nil {
		if cleanupErr := utils.DeleteFolder(folderPath); cleanupErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete folder: "+cleanupErr.Error())
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Conversion failed: "+err.Error())
		return
	}

	//upload the video file to cloudinary and return the file URL
	urlLink, err := utils.HandleFileUpload(ctx, videoPath, cfg.CloudinaryURL)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Remove the directory after uploading
	err = utils.DeleteFolder(folderPath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, ResponseMsg{Response: urlLink})
}
