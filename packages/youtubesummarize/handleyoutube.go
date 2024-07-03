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
// SOFTWARE
package youtubesummarize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/jpoz/groq"
)

type WhisperResponse struct {
	CompletedAt time.Time   `json:"completed_at"`
	CreatedAt   time.Time   `json:"created_at"`
	DataRemoved bool        `json:"data_removed"`
	Error       interface{} `json:"error"`
	ID          string      `json:"id"`
	Input       struct {
		URL          string `json:"url"`
		Task         string `json:"task"`
		Timestamp    string `json:"timestamp"`
		BatchSize    int    `json:"batch_size"`
		DiariseAudio bool   `json:"diarise_audio"`
	} `json:"input"`
	Logs    string `json:"logs"`
	Metrics struct {
		PredictTime float64 `json:"predict_time"`
		TotalTime   float64 `json:"total_time"`
	} `json:"metrics"`
	Output struct {
		Text   string `json:"text"`
		Chunks []struct {
			Text      string    `json:"text"`
			Timestamp []float64 `json:"timestamp"`
		} `json:"chunks"`
	} `json:"output"`
	StartedAt time.Time `json:"started_at"`
	Status    string    `json:"status"`
	Urls      struct {
		Get    string `json:"get"`
		Cancel string `json:"cancel"`
	} `json:"urls"`
	Version string `json:"version"`
}

const baseURL = "http://localhost:9000/api/v1/replicate/stt" //
// This function performs an api call to the replicate endpoint to transcribe the youtube video url and then sends the text to GPT AI MODEL to summarize.
func handleYoutubeSummariztion(youtubeURL, apiKey, grokApiKey string) (string, error) {
	// Create a buffer to write our form data to
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Create a form file field and write the url to it.
	if err := w.WriteField("url", youtubeURL); err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("model", "turian/insanely-fast-whisper-with-video")
	finalURL := baseURL + "?" + params.Encode()
	//create the http client
	// Close the multipart writer to set the terminating boundary
	if err := w.Close(); err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", finalURL, &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Api-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var WhisperResponseMessage WhisperResponse
	if resp.StatusCode != 200 {
		//handle error when status code is not 200
		if err := json.NewDecoder(resp.Body).Decode(&WhisperResponseMessage); err != nil {
			return "", fmt.Errorf("error decoding json response: %v", err)
		}
		return "", fmt.Errorf("error making api call to the whisper AI Model %v", WhisperResponseMessage.Error)
	}
	// The Api call succeded marshal the output and return the text.
	//Use the grok AI model to process the summarization.
	if err := json.NewDecoder(resp.Body).Decode(&WhisperResponseMessage); err != nil {
		return "", err
	}
	youtubeTranscribedText := WhisperResponseMessage.Output.Text

	grokClient := groq.NewClient(groq.WithAPIKey(grokApiKey))
	response, err := grokClient.CreateChatCompletion(groq.CompletionCreateParams{
		Model: "llama3-70b-8192",
		Messages: []groq.Message{
			{
				Role:    "system",
				Content: "You are a highly skilled AI model specialized in summarizing text transcribed From Youtube Videos. Your goal is to provide concise, accurate, and coherent summaries of the provided content.",
			},
			{
				Role:    "user",
				Content: youtubeTranscribedText,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to process summarization with grok model %v", err)
	}
	res := response.Choices[0].Message.Content
	return res, nil
}
