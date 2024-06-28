package musicdownloader

import (
	"encoding/json"
	"fmt"
	"github.com/kingmariano/omnicron-backendsever/golang-server/config"
	"github.com/kingmariano/omnicron-backendsever/golang-server/utils"
	"net/http"
)

type SongRequest struct {
	Song string `json:"song"`
}

// SongResponse represents the structure of the response
type SongResponse struct {
	Response []string `json:"response"`
}

var maxLength int64 //specifies the maxmium length of data returned from the youtube sdk
func DownloadMusic(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	ctx := r.Context()
	decode := json.NewDecoder(r.Body)
	params := SongRequest{}
	err := decode.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling json, %v", err))
		return
	}
	// create folder to handle downloads
	folderPath, err := utils.CreateUniqueFolder(utils.BasePath)
    if err!= nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
	//for accurate and precise result maxlength should be set to one.
	maxLength = 1
	audioDirectURL, err := downloadYoutubeLinkAndConvertToMp3(ctx, params.Song, maxLength, cfg.YoutubeDeveloperKey, cfg.CloudinaryUrl, folderPath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
    //clean up; remove folder after uploading
	err = utils.DeleteFolder(folderPath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, SongResponse{Response: audioDirectURL})
}
