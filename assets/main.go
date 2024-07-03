package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model,omitempty"`
	Stream   bool      `json:"stream,omitempty"`
	Proxy    string    `json:"proxy,omitempty"`
	Timeout  int       `json:"timeout,omitempty"`
	Shuffle  bool      `json:"shuffle,omitempty"`
}

type PythonResponse struct {
	Response interface{} `json:"response"`
}

func callPythonScript(requestBody RequestBody) (PythonResponse, error) {
	// Marshal the request body to JSON
	inputJSON, err := json.Marshal(requestBody)
	if err != nil {
		return PythonResponse{}, fmt.Errorf("failed to encode request body: %v", err)
	}

	// Execute the Python script
	cmd := exec.Command("python", "script.py")
	cmd.Stdin = bytes.NewReader(inputJSON)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return PythonResponse{}, fmt.Errorf("failed to execute script: %v, stderr: %s", err, stderr.String())
	}

	// Log stdout and stderr for debugging
	fmt.Printf("stdout: %s\n", out.String())
	fmt.Printf("stderr: %s\n", stderr.String())

	// Read and unmarshal the Python script response
	var pyResponse PythonResponse
	err = json.Unmarshal(out.Bytes(), &pyResponse)
	if err != nil {
		return PythonResponse{}, fmt.Errorf("failed to decode script response: %v", err)
	}

	return pyResponse, nil
}

func main() {
	// Example parameters to be sent to the Python script
	messages := []Message{
		{Role: "user", Content: "Hello"},
	}
	requestBody := RequestBody{
		Messages: messages,
		Model:    "",
		Stream:   false,
		Proxy:    "",
		Timeout:  10,
		Shuffle:  false,
	}

	// Call the Python script and get the response
	response, err := callPythonScript(requestBody)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the response
	fmt.Printf("Response: %+v\n", response)
}
