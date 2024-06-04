package generateimages

import (
	"context"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	rep "github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"net/http"
)

func ImageGeneration(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {

	ctx := context.Background()
	model := r.URL.Query().Get("model")
	if model == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "image model query parameter is required")
		return
	}
	repImageModel, err := rep.GetModelByName(model, rep.ImageModels)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "model not found")
		return
	}

	predictionFunc, ok := rep.Models[*repImageModel]
	if !ok {
		utils.RespondWithError(w, http.StatusNotFound, "model not found")
		return
	}
	modelIndx := rep.GetModelIndex(model, rep.ImageModels)

	predictionInput, err := processImageModelInput(repImageModel, ctx, r, modelIndx, cfg)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ImagePrediction, err := predictionFunc(ctx, cfg.ReplicateAPIKey, repImageModel.Version, predictionInput, nil, false)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, ImagePrediction)
}
