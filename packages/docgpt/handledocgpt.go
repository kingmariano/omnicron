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

	"github.com/h2non/filetype"
	"github.com/jpoz/groq"
	"github.com/kingmariano/omnicron/packages/gpt"
)

const baseURL = "http://localhost:8000/api/v1/doc_analyze" // URL to the doc analyze endpoint in the FastAPI server

// ErrorResponse represents the structure of error responses from the FastAPI server
type ErrorResponse struct {
	Detail string `json:"detail"`
}

// AnalyzeDocResponse represents the structure of the response from the "/doc_analyze" endpoint from the FastAPI server
type AnalyzeDocResponse struct {
	Text []string `json:"text"`
}

// CallDocGPTFastAPI calls the "/doc_analyze" endpoint from the FastAPI server and processes the response then uses the grok AI API client to make a request acting as a document gpt.
func CallDocGPTFastAPI(file multipart.File, fileHeader *multipart.FileHeader, prompt, grokAPIKey, apiKey, fastAPIBaseURL string) (string, error) {
	// Create a buffer to write our form data to
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Read the file into a byte slice
	filebytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Check if the file is empty
	if len(filebytes) == 0 {
		return "", errors.New("file is empty")
	}

	// Check file format if it's supported
	fileContentType, err := filetype.Get(filebytes)
	if err != nil {
		return "", fmt.Errorf("failed to get file type: %w", err)
	}
	if !checkContentType(fileContentType.Extension) {
		return "", fmt.Errorf("unsupported file format. supported formats are: %s", supportedFileTypes())
	}

	// Create a form file field
	fw, err := w.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy the file into the form field
	if _, err = io.Copy(fw, bytes.NewReader(filebytes)); err != nil {
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
	if docOutputText == "" {
		return "", errors.New("document analysis returned empty text")
	}
	docGptPrompt := docGPTPrompt(docOutputText)
	// Use the Groq library to create a chat completion request with the Groq API key
	grokClient := groq.NewClient(groq.WithAPIKey(grokAPIKey))
	response, err := grokClient.CreateChatCompletion(
		groq.CompletionCreateParams{
			Model: "mixtral-8x7b-32768",
			Messages: []groq.Message{
				{
					Role:    "system",
					Content: docGptPrompt,
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		log.Println("error occurred using the grok api; retrying with the g4f library")
		// if the grok AI model fails try using the g4f api Library.
		chatRequest := gpt.ChatRequest{
			Messages: []gpt.Message{
				{
					Role:    "system",
					Content: docGptPrompt,
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
		}
		g4fResponse, err := gpt.CallGPTFastAPI(chatRequest, apiKey, fastAPIBaseURL)
		if err != nil {
			return "", fmt.Errorf("failed to call gpt fast API: %w", err)
		}
		return g4fResponse.Response, nil
	}

	// Return the response from the Groq API
	return response.Choices[0].Message.Content, nil
}

func supportedFileTypes() []string {
	return []string{"pdf", "xps", "epub", "mobi", "fb2", "cbz", "svg", "txt"}
}

// function contains file types that are supported for document analyzing:
func checkContentType(fileExtension string) bool {
	supportedFileTypes := supportedFileTypes()
	for _, fileType := range supportedFileTypes {
		if fileExtension == fileType {
			return true
		}
	}
	return false
}

// input prompt to make the AI model behave as expected
func docGPTPrompt(documentText string) string {
	return fmt.Sprintf(`You are DocGPT, an advanced AI model specialized in interacting with and analyzing documents. Below is the text extracted from a document. Your task is to assist the user by responding to their queries about this document, providing clear and concise answers, summaries, and key points.

	Document Text: %s `, documentText)
}
