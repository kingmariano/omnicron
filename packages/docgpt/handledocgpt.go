package docgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/jpoz/groq"
)

const baseURL = "http://localhost:8000/api/v1/doc_analyze" // URL to the doc analyze endpoint in the Python server

// ErrorResponse represents the structure of error responses from the FastAPI server
type ErrorResponse struct {
	Detail string `json:"detail"`
}

// AnalyzeDocResponse represents the structure of the response from the document analysis
type AnalyzeDocResponse struct {
	Text []string `json:"text"`
}

// CallDocGPTFastAPI calls the "/doc_analyze" endpoint from the FastAPI server and processes the response then uses the grok AI API client to act as a document gpt.
func CallDocGPTFastAPI(file multipart.File, fileHeader *multipart.FileHeader, prompt, grokApiKey, apiKey string) (string, error) {
	// Create a buffer to write our form data to
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Create a form file field
	fw, err := w.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy the file into the form field
	if _, err = io.Copy(fw, file); err != nil {
		return "", fmt.Errorf("failed to copy file to form field: %w", err)
	}

	// Close the multipart writer to set the terminating boundary
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the HTTP client
	client := &http.Client{}

	// Create the HTTP request
	req, err := http.NewRequest("POST", baseURL, &b)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Api-Key", apiKey)

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		var errorMessage ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorMessage); err != nil {
			return "", fmt.Errorf("failed to decode error message: %w", err)
		}
		return "", errors.New(errorMessage.Detail)
	}

	// Decode the document analysis response
	var docResponse AnalyzeDocResponse
	if err := json.NewDecoder(resp.Body).Decode(&docResponse); err != nil {
		return "", fmt.Errorf("failed to decode document analysis response: %w", err)
	}
	log.Println("done analyzing document returning text")
	// Join the extracted text into a single string
	docOutputText := ""
	for _, text := range docResponse.Text {
		docOutputText += text + "\n"
	}

	// Use the Groq library to create a chat completion request with the Groq API key
	grokClient := groq.NewClient(groq.WithAPIKey(grokApiKey))
	response, err := grokClient.CreateChatCompletion(
		groq.CompletionCreateParams{
			Model: "mixtral-8x7b-32768",
			Messages: []groq.Message{
				{
					Role: "system",
					Content: fmt.Sprintf(`You are DocGPT, an advanced AI model specialized in interacting with and analyzing documents. Below is the text extracted from a document. Your task is to assist the user by responding to their queries about this document, providing clear and concise answers, summaries, and key points.

Document Text:
%s

Please respond to the following user query:`, docOutputText),
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion request: %w", err)
	}

	// Return the response from the Groq API
	return response.Choices[0].Message.Content, nil
}