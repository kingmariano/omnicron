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

package docgpt

import (
	"fmt"
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

type ResponseMsg struct {
	Response string `json:"response"`
}

func DocGPT(w http.ResponseWriter, r *http.Request, cfg *config.APIConfig) {
	// Parse the multipart form in the request
	err := r.ParseMultipartForm(30 << 20) // 30MB max memory
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing multipart form, %v", err))
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error retrieving the file, %v", err))
		return
	}
	defer file.Close()
	//prompt file value is required
	prompt := r.FormValue("prompt")
	if prompt == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Prompt is required")
		return
	}
	response, err := CallDocGPTFastAPI(file, fileHeader, prompt, cfg.GrokAPIKey, cfg.APIKey, cfg.FASTAPIBaseURL)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//responds with JSON
	utils.RespondWithJSON(w, http.StatusOK, ResponseMsg{Response: response})
}
