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

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// downloadFile downloads a file from the given URL and saves it to the specified path.

func DownloadFileURL(url, dest string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Add User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to download file: " + resp.Status)
	}

	// Extract file name and extension from URL
	fileName := path.Base(url)
	destPath := filepath.Clean(filepath.Join(dest, "/"+fileName))
	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Println(err)
		}
	}()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	// Normalize the file path format to use forward slashes
	normalizedPath := filepath.ToSlash(out.Name())

	return normalizedPath, nil
}

func ConvertFileToMP3(inputFile string) (string, error) {
	// Get the base name of the input file and change its extension to .mp3
	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	log.Print(baseName)
	outputFile := filepath.Clean(filepath.Join(filepath.Dir(inputFile), baseName+".mp3"))
	// Convert the input file to MP3 format
	err := ffmpeg.Input(inputFile).
		Output(outputFile, ffmpeg.KwArgs{"q:a": 0, "map": "a"}).
		Run()
	if err != nil {
		return "", fmt.Errorf("error converting file %s to mp3: %v", inputFile, err)
	}
	return outputFile, nil
}

// convertReaderToMP3 reads a video from an io.Reader and converts it to MP3.
func ConvertReaderToMP3(reader io.Reader, outputDir string) (string, error) {
	// Read the first 512 bytes to detect the file type
	buffer := make([]byte, 512)
	_, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read buffer: %w", err)
	}
	// Detect the file type based on the first 512 bytes
	contentType := http.DetectContentType(buffer)
	if !isValidContentType(contentType) {
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
	// Reset the reader to read the content again
	reader = io.MultiReader(bytes.NewReader(buffer), reader)
	// Create a temporary folder to store the video file and MP3 file

	// Create a temporary video file
	videoExtension := filepath.Ext(contentType)
	tempVideoFile := filepath.Clean(filepath.Join(outputDir, "input_video"+videoExtension))
	out, err := os.Create(tempVideoFile)
	if err != nil {
		return "", fmt.Errorf("failed to create video file: %w", err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Println(err)
		}
	}()
	// Copy the content of the reader to the video file
	_, err = io.Copy(out, reader)
	if err != nil {
		return "", fmt.Errorf("failed to copy content to video file: %w", err)
	}
	// Convert the video file to MP3
	mp3File, err := ConvertFileToMP3(tempVideoFile)
	if err != nil {
		return "", fmt.Errorf("failed to convert video file to MP3: %w", err)
	}
	return mp3File, nil
}

// isValidContentType checks if the content type is valid for conversion.
func isValidContentType(contentType string) bool {
	prefix := strings.SplitN(contentType, "/", 2)[0]
	return prefix == "video" || prefix == "audio"
}
