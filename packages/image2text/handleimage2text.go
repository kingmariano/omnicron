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
package image2text

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8000/api/v1/image_to_text" //url to the image to text endpoint in the python server
// Calls the "/image+to_text" endpoint from the fastAPI server
type ImageToTextResponse struct {
	Text string `json:"text"`
}
type ErrorResponse struct {
	Detail string `json:"detail"`
}

func CallImageToTextFastAPI(file multipart.File, fileHeader *multipart.FileHeader, apiKey string) (*ImageToTextResponse, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Create a form file field
	fw, err := w.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, err
	}
	// Copy the file into the form field
	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}
	// Close the multipart writer to set the terminating boundary
	if err := w.Close(); err != nil {
		return nil, err
	}
	//create the http client
	client := &http.Client{
		Timeout: time.Second * 300,
	}
	req, err := http.NewRequest("POST", baseURL, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Api-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		//handle error when status code is not 200
		var errorMessage ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorMessage); err != nil {
			return nil, err
		}
		return nil, errors.New(errorMessage.Detail)
	}
	var response ImageToTextResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil

}
