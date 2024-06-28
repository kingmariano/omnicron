package convert2mp3

import (
	"net/http"

	"github.com/kingmariano/omnicron-backendsever/golang-server/config"
	"github.com/kingmariano/omnicron-backendsever/golang-server/utils"
)

type Response struct {
	Response string `json:"response"`
}

func ConvertToMp3(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	ctx := r.Context()
	folderPath, err := utils.CreateUniqueFolder(utils.BasePath)
	if err!= nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
	outputfileName, err := handleRequestBodyAndConvertToMP3(r, folderPath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// upload file to cloudinary
	urlLink, err := utils.HandleFileUpload(ctx, outputfileName, cfg.CloudinaryUrl)
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
	utils.RespondWithJSON(w, http.StatusOK, Response{Response: urlLink})
}
