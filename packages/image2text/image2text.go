package image2text

import (
	"fmt"
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

type ImageToTextResponse struct {
	Text string `json:"text"`
}
type ErrorResponse struct {
	Detail string `json:"detail"`
}

func Image2text(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	// Parse the multipart form in the request
	err := r.ParseMultipartForm(30 << 20) // 30MB max memory
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing multipart form, %v", err))
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error retrieving the file, %v", err))
		return
	}
	defer file.Close()
	response, err := CallImageToTextFastAPI(file, fileHeader, cfg.ApiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error calling the Image To Text Endpoint, %v", err))
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}
