package youtubesummarize

import (
	"encoding/json"
	"fmt"
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

type YoutubeRequest struct {
	URL string `json:"url"`
}

func YoutubeSummarization(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	decode := json.NewDecoder(r.Body)
	youtubeParams := YoutubeRequest{}
	err := decode.Decode(&youtubeParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling json, %v", err))
		return
	}
	summary, err := handleYoutubeSummariztion(youtubeParams.URL, cfg.ApiKey, cfg.GrokApiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, summary)
}
