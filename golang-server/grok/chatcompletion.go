package grok

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"github.com/jpoz/groq"
)

func ChatCompletion(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	decode := json.NewDecoder(r.Body)
	grokParams := groq.CompletionCreateParams{}
	err := decode.Decode(&grokParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling json, %v", err))
		return
	}

	client := groq.NewClient(groq.WithAPIKey(cfg.GrokApiKey))
	response, err := client.CreateChatCompletion(grokParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error handling chat completion, %v", err))
		return
	}
	res := response.Choices[0].Message.Content
	utils.RespondWithJSON(w, http.StatusOK, res)
}
