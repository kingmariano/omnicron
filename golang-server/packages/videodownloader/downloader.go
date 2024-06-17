package videodownloader

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	"github.com/charlesozo/omnicron-backendsever/golang-server/storage"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"net/http"
	"os"
	"path/filepath"
)

type DownloadParams struct {
	URL        string `json:"url"`
	Resolution string `json:"resolution"`
}
type Responseparams struct {
	Response string `json:"response"`
}

func DownloadVideo(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	ctx := context.Background()

	decode := json.NewDecoder(r.Body)
	params := DownloadParams{}
	err := decode.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling json, %v", err))
		return
	}
	outputName := "youtube"
	outputPath := "./downloadedvideo/"
	// Check if the folder "downloadedvideo" exists, create if not
	// if _, err := os.Stat(outputPath); os.IsNotExist(err) {
	// 	err := os.Mkdir(outputPath, os.ModePerm)
	// 	if err != nil {
	// 		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating directory, %v", err))
	// 		return
	// 	}
	// 	defer os.RemoveAll(outputPath) // Ensure the folder is deleted after task completion
	// }
	 utils.CreateFolder(outputPath)
	// if err != nil{
	// 	utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating directory, %v", err))
	// 	return
	// }

	err = DownloadVideoData(params.URL, outputName, outputPath, params.Resolution)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	files, err := filepath.Glob(outputPath + outputName + ".*")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(files) == 0 {
		utils.RespondWithError(w, http.StatusInternalServerError, "No video file found")
		return
	}
	videoPath := files[0]
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting file info, %v", err))
		return
	}
	if fileInfo.Size() == 0 {
		utils.RespondWithError(w, http.StatusInternalServerError, "Downloaded video file is empty")
		return
	}
	urlLink, err := storage.HandleFileUpload(ctx, videoPath, cfg)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Remove the video file after uploading
	err = os.Remove(videoPath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, Responseparams{Response: urlLink})
}
