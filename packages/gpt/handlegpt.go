package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8000/api/v1/chat/completion" //url to the chat completion endpoint in the python server
// Calls the "/chatcompletion" endpoint from the fastAPI server
func CallGPTFastAPI(request ChatRequest, apiKey string) (*ChatResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		//handle error when status code is not 200
		var errorMessage ChatError
		if err := json.NewDecoder(resp.Body).Decode(&errorMessage); err != nil {
			return nil, err
		}
		return nil, errors.New(errorMessage.Detail)
	}
	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
