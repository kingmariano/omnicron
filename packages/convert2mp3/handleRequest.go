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

package convert2mp3

import (
	"errors"
	"fmt"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

func handleRequestBodyAndConvertToMP3(r *http.Request, outputDir string) (string, error) {
	// handles the request body input which is either a file-form or a url form-value
	url := r.FormValue("url")
	//if the user has specified the url parameter handles it immedaitely
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

	// if the url parameter is not set but a multipart file form is specified
	file, _, err := r.FormFile("file")
	if err != nil {
		// If there's an error and no URL was provided, return the error
		return "", errors.New("please provide either a valid URL or a file")
	}
	defer file.Close()
	// performs the conversion of the reader to mp3
	outputFileName, err := utils.ConvertReaderToMP3(file, outputDir)
	if err != nil {
		return "", fmt.Errorf("error converting uploaded file to mp3: %v", err)
	}
	return outputFileName, nil
}
