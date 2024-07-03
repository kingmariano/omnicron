package gpt

import (
	"encoding/json"
	"fmt"
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

type ChatRequest struct {
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Model    string `json:"model,omitempty"`
	Stream   bool   `json:"stream,omitempty"`
	Proxy    string `json:"proxy,omitempty"`
	Timeout  int    `json:"timeout,omitempty"`
	Shuffle  bool   `json:"shuffle,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}
type ChatResponse struct {
	Response string `json:"response"`
}
type ErrorResponse struct {
	Detail string `json:"detail"`
}

func ChatCompletion(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	decode := json.NewDecoder(r.Body)
	chatParams := ChatRequest{}
	err := decode.Decode(&chatParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling json, %v", err))
		return
	}
	response, err := CallGPTFastAPI(chatParams, cfg.ApiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error handling chat completion, %v", err))
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, response)

}
