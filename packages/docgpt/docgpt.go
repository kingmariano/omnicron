package docgpt

import (
	"fmt"
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

func DocGPT(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
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
	prompt := r.FormValue("prompt")
	if prompt == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Prompt is required")
		return
	}
	response, err := CallDocGPTFastAPI(file, fileHeader, prompt, cfg.GrokApiKey, cfg.ApiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}
