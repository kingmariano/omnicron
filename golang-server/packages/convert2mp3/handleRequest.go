package convert2mp3

import (
	"errors"
	"fmt"
	"net/http"
	"github.com/kingmariano/omnicron-backendsever/golang-server/utils"
)

func handleRequestBodyAndConvertToMP3(r *http.Request, outputDir string) (string, error) {
	// Handle URL input
	url := r.FormValue("url")
	if url != "" {
		downloadedFileName, err := utils.DownloadFileURL(url, outputDir)
		if err != nil {
			return "", fmt.Errorf("error downloading file %s: %v", url, err)
		}
		outputFileName, err := utils.ConvertFileToMP3(downloadedFileName)
		if err != nil {
			return "", fmt.Errorf("error converting file %s to mp3: %v", downloadedFileName, err)
		}
		return outputFileName, nil
	}

	// Handle file input
	file, _, err := r.FormFile("file")
	if err != nil {
		// If there's an error and no URL was provided, return the error
		return "", errors.New("please provide either a valid URL or a file")
	}
	defer file.Close()

	outputFileName, err := utils.ConvertReaderToMP3(file, outputDir)
	if err != nil {
		return "", fmt.Errorf("error converting uploaded file to mp3: %v", err)
	}
	return outputFileName, nil
}
