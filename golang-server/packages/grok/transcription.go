package grok

import (
	"fmt"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"github.com/jpoz/groq"
	"io"
	"net/http"
	"os"
)

func Transcription(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // limit your max input length
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing multipart form, %v", err))
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error retrieving the file, %v", err))
		return
	}
	defer file.Close()
	// Create temporary file to store uploaded file
	tempFile, err := os.CreateTemp("", handler.Filename)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating temp file, %v", err))
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	// Copy the file data to the temp file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error copying file data to temp file: %v", err))
		return
	}

	// Re-open the file as *os.File
	osFile, err := os.Open(tempFile.Name())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error re-opening temp file: %v", err))
		return
	}
	defer osFile.Close()

	// Now you can use osFile as *os.File
	// For demonstration purposes, we'll just print the file name
	fmt.Printf("File uploaded successfully: %s\n", tempFile.Name())
	model := r.FormValue("model")
	language := r.FormValue("language")

	grokParams := groq.TranscriptionCreateParams{
		File:     osFile,
		Model:    groq.TranslationModel(model), // You may need to convert this to TranslationModel
		Language: language,
	}
	client := groq.NewClient(groq.WithAPIKey(cfg.GrokApiKey))
	response, err := client.CreateTranscription(grokParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error Transcribing: %v", err))
		return
	}
	res := response.Choices[0].Message.Content
	utils.RespondWithJSON(w, http.StatusOK, res)
}
