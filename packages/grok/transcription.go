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
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jpoz/groq" // Import groq package for transcription
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
)

// Transcription handles HTTP requests to perform transcription.
func Transcription(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {
	// Parse the multipart form with a max size of 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing multipart form: %v", err))
		return
	}

	// Retrieve the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error retrieving the file: %v", err))
		return
	}
	defer file.Close()

	// Create a temporary file to store the uploaded file
	tempFile, err := os.CreateTemp("", handler.Filename)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating temp file: %v", err))
		return
	}
	defer os.Remove(tempFile.Name()) // Clean up: remove the temporary file
	defer tempFile.Close()

	// Copy the file data to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error copying file data to temp file: %v", err))
		return
	}

	// Re-open the temporary file as *os.File
	osFile, err := os.Open(tempFile.Name())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error re-opening temp file: %v", err))
		return
	}
	defer osFile.Close()

	// Prepare transcription parameters
	model := r.FormValue("model")
	language := r.FormValue("language")
	grokParams := groq.TranscriptionCreateParams{
		File:     osFile,
		Model:    groq.TranslationModel(model), // Ensure correct conversion to TranslationModel type
		Language: language,
	}

	// Create a groq client with API key from config
	client := groq.NewClient(groq.WithAPIKey(cfg.GrokAPIKey))
	response, err := client.CreateTranscription(grokParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error transcribing: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.ResponseMsg{Response: response})
}
