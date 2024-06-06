package generatevideos

import (
	"context"
	"errors"
	"fmt"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	rep "github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate"
	"github.com/charlesozo/omnicron-backendsever/golang-server/storage"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	replicate "github.com/replicate/replicate-go"
	"log"
	"net/http"
)

func processVideoModelInput(imageModel *rep.ReplicateModel, ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	if imageModel.Category == "High" {
		replicateInput, err := processHighVideoInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	}
	return nil, errors.New("image category unavailable")
}

func processHighVideoInput(ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	log.Println("This is  High video generation")
	var HighVideoGenerationParams rep.HighVideoGenerationParams
	err := r.ParseMultipartForm(50 << 20) // 50MB
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}

	HighVideoGenerationParams = rep.HighVideoGenerationParams{}.ZeroscopeV2Xl()
	log.Println(HighVideoGenerationParams)
	prompt := r.FormValue("prompt")
	if prompt == "" {
		return nil, errors.New("prompt cant be empty")
	}
	HighVideoGenerationParams.Prompt = prompt
	//extract and replace default
	utils.SetStringValue(r.FormValue("negative_prompt"), &HighVideoGenerationParams.NegativePrompt)
	utils.SetStringValue(r.FormValue("init_video_url"), &HighVideoGenerationParams.InitVideoURL)
	utils.SetFloatValue(r.FormValue("init_weight"), &HighVideoGenerationParams.InitWeight)
	utils.SetIntValue(r.FormValue("num_frames"), &HighVideoGenerationParams.NumFrames)
	utils.SetIntValue(r.FormValue("num_inferences_steps"), &HighVideoGenerationParams.NumInferenceSteps)
	utils.SetIntValue(r.FormValue("width"), &HighVideoGenerationParams.Width)
	utils.SetIntValue(r.FormValue("height"), &HighVideoGenerationParams.Height)
	utils.SetFloatValue(r.FormValue("guidance_scale"), &HighVideoGenerationParams.GuidanceScale)
	utils.SetIntValue(r.FormValue("fps"), &HighVideoGenerationParams.FPS)
	utils.SetStringValue(r.FormValue("video_model"), &HighVideoGenerationParams.VideoModel)
	utils.SetIntValue(r.FormValue("batch_size"), &HighVideoGenerationParams.BatchSize)
	utils.SetBoolValue(r.FormValue("remove_watermark"), &HighVideoGenerationParams.RemoveWatermark)
	// Handle initial video file
	videoFile, _, err := r.FormFile("video")
	if err == nil {
		videoUrl, err := storage.HandleFileUpload(ctx, videoFile, cfg)
		if err != nil {
			return nil, err
		}
		HighVideoGenerationParams.InitVideoURL = &videoUrl
	}
	if videoFile != nil {
		defer videoFile.Close()
	}

	input := replicate.PredictionInput{
		"prompt":              HighVideoGenerationParams.Prompt,
		"fps":                 HighVideoGenerationParams.FPS,
		"model":               HighVideoGenerationParams.VideoModel,
		"width":               HighVideoGenerationParams.Width,
		"height":              HighVideoGenerationParams.Height,
		"batch_size":          HighVideoGenerationParams.BatchSize,
		"num_frames":          HighVideoGenerationParams.NumFrames,
		"init_weight":         HighVideoGenerationParams.InitWeight,
		"guidance_scale":      HighVideoGenerationParams.GuidanceScale,
		"negative_prompt":     HighVideoGenerationParams.NegativePrompt,
		"remove_watermark":    HighVideoGenerationParams.RemoveWatermark,
		"num_inference_steps": HighVideoGenerationParams.NumInferenceSteps,
	}
	if HighVideoGenerationParams.InitVideoURL != nil {
		input["init_video"] = *HighVideoGenerationParams.InitVideoURL
	}
	return input, nil

}
