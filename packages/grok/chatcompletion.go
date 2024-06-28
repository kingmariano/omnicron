package grok

import (
	"encoding/json"
	"fmt"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/jpoz/groq"
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

func validateParams(g groq.CompletionCreateParams) error {
	return validation.ValidateStruct(&g, validation.Field(&g.Model, validation.Required, validation.In("llama3-8b-8192", "llama3-70b-8192", "mixtral-8x7b-32768", " gemma-7b-it")),
		validation.Field(&g.Messages, validation.Required),
	)

}
func ChatCompletion(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	decode := json.NewDecoder(r.Body)
	grokParams := groq.CompletionCreateParams{}
	err := decode.Decode(&grokParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling json, %v", err))
		return
	}
	// Validate the request parameters
	err = validateParams(grokParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error, %v", err))
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
