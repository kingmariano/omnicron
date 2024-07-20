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
	"net/http"

	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
)

func ConvertToMp3(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {
	ctx := r.Context()
	// creates a unique folder within the current directory
	folderPath, err := utils.CreateUniqueFolder(utils.BasePath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// processes the uploaded file and converts it to mp3, then saves it to the unique folder path
	outputfileName, err := handleRequestBodyAndConvertToMP3(r, folderPath)
	if err != nil {
		if cleanupErr := utils.DeleteFolder(folderPath); cleanupErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete folder: "+cleanupErr.Error())
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Conversion failed: "+err.Error())
		return
	}
	// uploads the file to cloudinary to get back the direct url link
	urlLink, err := utils.HandleFileUpload(ctx, outputfileName, cfg.CloudinaryURL)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Remove the directory after uploading
	err = utils.DeleteFolder(folderPath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//returns back the response in JSON format
	utils.RespondWithJSON(w, http.StatusOK, utils.ResponseMsg{Response: urlLink})
}
