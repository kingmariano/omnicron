package stt

import (
	"context"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	rep "github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"net/http"
)

func STT(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	ctx := context.Background()
	model := r.URL.Query().Get("model")
	if model == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "stt model query parameter is required")
		return
	}
	repSTTModel, err := rep.GetModelByName(model, rep.STTModels)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "model not found")
		return
	}
	predictionFunc, ok := rep.STTGenModels[*repSTTModel]
	if !ok {
		utils.RespondWithError(w, http.StatusNotFound, "model not found")
		return
	}
	predictionInput, err := processSTTModelInput(repSTTModel, ctx, r, cfg)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	STTPrediction, err := predictionFunc(ctx, cfg.ReplicateAPIKey, repSTTModel.Version, predictionInput, nil, false)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, STTPrediction)

}
