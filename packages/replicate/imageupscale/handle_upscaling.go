package imageupscale

import (
	"context"
	"errors"
	"fmt"
	"github.com/kingmariano/omnicron/config"
	rep "github.com/kingmariano/omnicron/packages/replicate"
	"github.com/kingmariano/omnicron/utils"
	replicate "github.com/replicate/replicate-go"
	"log"
	"net/http"
)

func processImageUpscaleModelInput(imageModel *rep.ReplicateModel, ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	if imageModel.Category == "High" {
		replicateInput, err := processHighUpscalingInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	} else if imageModel.Category == "Low" {
		replicateInput, err := processLowUpscalingInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	}
	return nil, errors.New("image category unavailable")
}
func processLowUpscalingInput(ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	log.Println("This is  low upscaling")
	var LowImageUpscaleGenerationParams rep.LowImageUpscaleGenerationParams
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}
	//set the parameters to their default
	LowImageUpscaleGenerationParams = rep.LowImageUpscaleGenerationParams{}.RealEsrgan()
	_, imageFileHeader, err := r.FormFile("image")
	if err != nil {
		return nil, fmt.Errorf("provide image file: %v", err)
	}
	repFile, err := rep.RequestFileToReplicateFile(ctx, imageFileHeader, cfg.ReplicateAPIKey)
	if err != nil {
		return nil, err
	}
	LowImageUpscaleGenerationParams.ImageFile = repFile
	utils.SetFloatValue(r.FormValue("scale"), &LowImageUpscaleGenerationParams.Scale)
	utils.SetBoolValue(r.FormValue("face_enhance"), &LowImageUpscaleGenerationParams.FaceEnhance)
	input := replicate.PredictionInput{
		"image":        LowImageUpscaleGenerationParams.ImageFile,
		"scale":        LowImageUpscaleGenerationParams.Scale,
		"face_enhance": LowImageUpscaleGenerationParams.FaceEnhance,
	}
	return input, nil
}

func processHighUpscalingInput(ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	log.Println("This is  high upscaling")
	var HighImageUpscaleGenerationParams rep.HighImageUpscaleGenerationParams
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}
	//set the parameters to their default
	HighImageUpscaleGenerationParams = rep.HighImageUpscaleGenerationParams{}.NewClarityUpscaler()

	_, imageFileHeader, err := r.FormFile("image")
	if err != nil {
		return nil, fmt.Errorf("provide image file: %v", err)
	}
	repFile, err := rep.RequestFileToReplicateFile(ctx, imageFileHeader, cfg.ReplicateAPIKey)
	if err != nil {
		return nil, err
	}
	HighImageUpscaleGenerationParams.ImageFile = repFile
	// Extract form values and set them
	utils.SetStringValue(r.FormValue("prompt"), &HighImageUpscaleGenerationParams.Prompt)
	utils.SetStringValue(r.FormValue("negative_prompt"), &HighImageUpscaleGenerationParams.NegativePrompt)
	utils.SetFloatValue(r.FormValue("scale_factor"), &HighImageUpscaleGenerationParams.ScaleFactor)
	utils.SetFloatValue(r.FormValue("dynamic"), &HighImageUpscaleGenerationParams.Dynamic)
	utils.SetFloatValue(r.FormValue("creativity"), &HighImageUpscaleGenerationParams.Creativity)
	utils.SetFloatValue(r.FormValue("resemblance"), &HighImageUpscaleGenerationParams.Resemblance)
	utils.SetIntValue(r.FormValue("tiling_width"), &HighImageUpscaleGenerationParams.TilingWidth)
	utils.SetIntValue(r.FormValue("tiling_height"), &HighImageUpscaleGenerationParams.TilingHeight)
	utils.SetStringValue(r.FormValue("sd_model"), &HighImageUpscaleGenerationParams.SdModel)
	utils.SetStringValue(r.FormValue("scheduler"), &HighImageUpscaleGenerationParams.Scheduler)
	utils.SetIntValue(r.FormValue("num_inference_steps"), &HighImageUpscaleGenerationParams.NumInferenceSteps)
	utils.SetIntValue(r.FormValue("seed"), &HighImageUpscaleGenerationParams.Seed)
	utils.SetBoolValue(r.FormValue("downscaling"), &HighImageUpscaleGenerationParams.Downscaling)
	utils.SetIntValue(r.FormValue("downscaling_resolution"), &HighImageUpscaleGenerationParams.DownscalingResolution)
	utils.SetFloatValue(r.FormValue("sharpen"), &HighImageUpscaleGenerationParams.Sharpen)
	utils.SetStringValue(r.FormValue("output_format"), &HighImageUpscaleGenerationParams.OutputFormat)
	input := replicate.PredictionInput{
		"image":                  HighImageUpscaleGenerationParams.ImageFile,
		"prompt":                 HighImageUpscaleGenerationParams.Prompt,
		"negative_prompt":        HighImageUpscaleGenerationParams.NegativePrompt,
		"scale_factor":           HighImageUpscaleGenerationParams.ScaleFactor,
		"seed":                   HighImageUpscaleGenerationParams.Seed,
		"dynamic":                HighImageUpscaleGenerationParams.Dynamic,
		"sharpen":                HighImageUpscaleGenerationParams.Sharpen,
		"sd_model":               HighImageUpscaleGenerationParams.SdModel,
		"scheduler":              HighImageUpscaleGenerationParams.Scheduler,
		"creativity":             HighImageUpscaleGenerationParams.Creativity,
		"downscaling":            HighImageUpscaleGenerationParams.Downscaling,
		"resemblance":            HighImageUpscaleGenerationParams.Resemblance,
		"tiling_width":           HighImageUpscaleGenerationParams.TilingWidth,
		"tiling_height":          HighImageUpscaleGenerationParams.TilingHeight,
		"num_inference_steps":    HighImageUpscaleGenerationParams.NumInferenceSteps,
		"output_format":          HighImageUpscaleGenerationParams.OutputFormat,
		"downscaling_resolution": HighImageUpscaleGenerationParams.DownscalingResolution,
	}

	return input, nil
}
