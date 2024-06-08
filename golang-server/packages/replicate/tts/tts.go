package tts
import (
	"context"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	rep "github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"net/http"
)
func TTS(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig){
	ctx := context.Background()
	model := r.URL.Query().Get("model")
	if model == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "tts model query parameter is required")
		return
	}
	repTTSModel, err := rep.GetModelByName(model, rep.TTSModels)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "model not found")
		return
	}
	predictionFunc, ok := rep.TTSGenModels[*repTTSModel]
	if !ok {
		utils.RespondWithError(w, http.StatusNotFound, "model not found")
		return
	}
	predictionInput, err := processTTSModelInput(repTTSModel, ctx, r, cfg)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	TTSPrediction, err := predictionFunc(ctx, cfg.ReplicateAPIKey, repTTSModel.Version, predictionInput, nil, false)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, TTSPrediction)
	
}