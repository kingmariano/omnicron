package shazam

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8000/api/v1/shazam" //url to the shazam endpoint in the python server
// Calls the "/shazam" endpoint from the fastAPI server
func CallShazamFastAPI(file multipart.File, fileHeader *multipart.FileHeader, apiKey string) (*ShazamResponse, error) {
	// Create a buffer to write our form data to
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
	var response ShazamResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
