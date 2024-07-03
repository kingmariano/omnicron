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

package utils

import (
	"fmt"
	"github.com/google/uuid"
	"os"
)

// set the default path to where the video or audio will be downloaded
var BasePath string = "./downloads" // Base path for downloads
var OutputName string = "youtube"   // Default output name

// CreateUniqueFolder creates a new folder with a unique name within the specified base path
func CreateUniqueFolder(basePath string) (string, error) {
	uniqueFolder := basePath + uuid.New().String() // Generate a unique folder name by appending a UUID to the base path
	// Create the directory with the specified permissions (0750)
	err := os.MkdirAll(uniqueFolder, 0750)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}
	// Return the unique folder path if successful
	return uniqueFolder, nil
}

// DeleteFolder removes the specified folder and all its contents
func DeleteFolder(path string) error {
	// Remove the directory and all its contents
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("error deleting directory: %v", err)
	}
	// Return nil if deletion is successful
	return nil
}
