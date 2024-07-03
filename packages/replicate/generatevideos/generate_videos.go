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
// SOFTWARE
package generatevideos

import (
	"github.com/kingmariano/omnicron/config"
	rep "github.com/kingmariano/omnicron/packages/replicate"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

func VideoGeneration(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	ctx := r.Context()
	model := r.URL.Query().Get("model")
	if model == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "video model query parameter is required")
		return
	}
	repVideoModel, err := rep.GetModelByName(model, rep.VideoModels)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "model not found")
		return
	}
	predictionFunc, ok := rep.VideoGenModels[*repVideoModel]
	if !ok {
		utils.RespondWithError(w, http.StatusNotFound, "model not found")
		return
	}
	predictionInput, err := processVideoModelInput(repVideoModel, ctx, r, cfg)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	videoPrediction, err := predictionFunc(ctx, cfg.ReplicateAPIKey, repVideoModel.Version, predictionInput, nil, false)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, videoPrediction)
}
