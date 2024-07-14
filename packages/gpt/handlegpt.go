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
package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Calls the "/chatcompletion" endpoint from the fastAPI server
func CallGPTFastAPI(request ChatRequest, apiKey string, fastAPIBaseURL string) (*ChatResponse, error) {
	fastAPIGPTEndPoint := fmt.Sprintf("%s/api/v1/chat/completion", fastAPIBaseURL) //url to the chat completion endpoint in the FastAPI server
	// marshal the request to json format for sending to the server
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest("POST", fastAPIGPTEndPoint, bytes.NewBuffer(jsonData))
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
	//handle error when status code is not 200
	if resp.StatusCode != 200 {
		var errorMessage ErrorResponse
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
