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
package grok

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-ozzo/ozzo-validation" // Import validation package for input validation
	"github.com/jpoz/groq"               // Import groq package for chat completions
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
)

// validateParams validates the input parameters for creating a chat completion.
func validateParams(g groq.CompletionCreateParams) error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.Model, validation.Required, validation.In("llama3-8b-8192", "llama3-70b-8192", "mixtral-8x7b-32768", "gemma-7b-it")), // Validate the 'Model' field
		validation.Field(&g.Messages, validation.Required), // Validate the 'Messages' field
	)
}

// ChatCompletion handles HTTP requests to create a chat completion.
func ChatCompletion(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {
	decode := json.NewDecoder(r.Body) // Create a JSON decoder for decoding request body
	grokParams := groq.CompletionCreateParams{}
	err := decode.Decode(&grokParams) // Decode JSON request body into grokParams struct
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

	// Create a new groq client with API key from config
	client := groq.NewClient(groq.WithAPIKey(cfg.GrokAPIKey))
	response, err := client.CreateChatCompletion(grokParams) // Call groq API to create chat completion
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error handling chat completion, %v", err))
		return
	}

	// Get the response message content from the first choice

	// new release for grok response
	// Respond with JSON containing the response message
	utils.RespondWithJSON(w, http.StatusOK,  response)
}
