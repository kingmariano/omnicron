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
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	ffmpeg "github.com/u2takey/ffmpeg-go" // Import the ffmpeg-go package for video and audio processing
)
var UnsupportedFileFormat error = errors.New("unsupported file format")
// DownloadFileURL downloads a file from the given URL and saves it to the specified path.
func DownloadFileURL(url, dest string) (string, error) {
	client := &http.Client{}                     // Create an HTTP client
	req, err := http.NewRequest("GET", url, nil) // Create a new GET request
	if err != nil {
		return "", err // Return an error if request creation fails
	}

	// Add User-Agent header to the request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to download file: " + resp.Status) // Return an error if the status is not OK
	}

	// Extract file name and extension from URL
	fileName := path.Base(url)
	destPath := filepath.Clean(filepath.Join(dest, "/"+fileName)) // Clean and join the destination path and file name
	// Create the destination file
	out, err := os.Create(destPath)
	if err != nil {
		return "", err // Return an error if file creation fails
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Println(err) // Log an error if closing the file fails
		}
	}()

	_, err = io.Copy(out, resp.Body) // Copy the response body to the destination file
	if err != nil {
		return "", err // Return an error if copying fails
	}

	// Normalize the file path format to use forward slashes
	normalizedPath := filepath.ToSlash(out.Name())

	return normalizedPath, nil // Return the normalized file path
}

// ConvertFileToMP3 converts the given input filepath to MP3 format.
func ConvertFileToMP3(inputFilePath string) (string, error) {
	// Get the base name of the input file and change its extension to .mp3
	baseName := strings.TrimSuffix(filepath.Base(inputFilePath), filepath.Ext(inputFilePath))
	log.Print(baseName)
	outputFile := filepath.Clean(filepath.Join(filepath.Dir(inputFilePath), baseName+".mp3")) // Clean and join the output path and file name

	// Convert the input file to MP3 format using ffmpeg
	err := ffmpeg.Input(inputFilePath).
		Output(outputFile, ffmpeg.KwArgs{"q:a": 0, "map": "a"}). // Set output options for MP3
		Run()
	if err != nil {
		return "", fmt.Errorf("error converting file %s to mp3: %v", inputFilePath, err) // Return an error if conversion fails
	}
	return outputFile, nil // Return the output file path
}

// ConvertReaderToMP3 reads a video from an io.Reader and converts it to MP3.
func ConvertReaderToMP3(reader io.Reader, outputDir string) (string, error) {
   // read the content of the file being given convert to bytes to detect file type
   fileBytes, err := io.ReadAll(reader)
   if err != nil {
	    return "", fmt.Errorf("failed to read file content: %w", err) // Return an error if reading fails
   }

   // Detect the file type based on the content type
    contentType, err :=  getSupportedFileType(fileBytes)
	if err!= nil {
        return "", err
    }
	
	videoExtension := filepath.Ext(contentType)
	tempVideoFile := filepath.Clean(filepath.Join(outputDir, "input_video"+videoExtension)) // Clean and join the temp video file path
	out, err := os.Create(tempVideoFile) // Create the temporary video file
	if err != nil {
		return "", fmt.Errorf("failed to create video file: %w", err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Println(err) // Log an error if closing the file fails
		}
	}()
	// Create a new reader from the fileBytes
	newReader := bytes.NewReader(fileBytes)

	// Copy the content of the reader to the video file
	_, err = io.Copy(out, newReader)
	if err != nil {
		return "", fmt.Errorf("failed to copy content to video file: %w", err)
	}

	// Convert the video file to MP3
	mp3File, err := ConvertFileToMP3(tempVideoFile)
	if err != nil {
		return "", fmt.Errorf("failed to convert video file to MP3: %w", err) // Return an error if conversion fails
	}

	return mp3File, nil // Return the MP3 file path
}
//supported file format for conversion to Mp3
func supportedFileFormats()[]string{
  return []string{"video", "audio"}
}
// retrieves the content type of the file if the file is supported
func getSupportedFileType(fileType []byte) (string, error){
	fileFormat, err := filetype.Get(fileType)
	if err!= nil {
        return "", fmt.Errorf("failed to get file type: %w", err)
    }
	supportedFileFormat := supportedFileFormats()
	for _, fileType := range supportedFileFormat{
       if fileFormat.MIME.Type == fileType{
		return fileFormat.MIME.Value, nil
	   }
	}
    return "", UnsupportedFileFormat
}
