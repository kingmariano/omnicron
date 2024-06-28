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
