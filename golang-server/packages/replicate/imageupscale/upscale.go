package imageupscale

import (
	"context"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	rep "github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"net/http"
)

func ImageUpscale(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	ctx := context.Background()
	model := r.URL.Query().Get("model")
	if model == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "image model query parameter is required")
		return
	}
	repImageUpscaleModel, err := rep.GetModelByName(model, rep.ImageUpscaleModels)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "model not found")
		return
	}
	predictionFunc, ok := rep.ImageUpscaleGenModels[*repImageUpscaleModel]
	if !ok {
		utils.RespondWithError(w, http.StatusNotFound, "model not found")
		return
	}
	predictionInput, err := processImageUpscaleModelInput(repImageUpscaleModel, ctx, r, cfg)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ImageUpscalePrediction, err := predictionFunc(ctx, cfg.ReplicateAPIKey, repImageUpscaleModel.Version, predictionInput, nil, false)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, ImageUpscalePrediction)
}
