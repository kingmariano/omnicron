package generateimages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kingmariano/omnicron/config"
	rep "github.com/kingmariano/omnicron/packages/replicate"
	"github.com/kingmariano/omnicron/utils"
	replicate "github.com/replicate/replicate-go"
	"net/http"
)

//	var imagemodelInputProcessors = map[int]func(ctx context.Context, r *http.Request, token string, modelIndex int) (replicate.PredictionInput, error){
//		0: processLowImageGenerationInput,
//		1: processHighImageGenerationInput,
//		2: processHighImageGenerationInput,
//		3: processLowImageGenerationInput,
//		4: processHighImageGenerationInput,
//		// Add other model processors similarly...
//	}

func processImageModelInput(imageModel *rep.ReplicateModel, ctx context.Context, r *http.Request, modelIndex int, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	if imageModel.Category == "Low" {
		replicateInput, err := processLowImageGenerationInput(ctx, r, modelIndex, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	} else if imageModel.Category == "High" {
		replicateInput, err := processHighImageGenerationInput(ctx, r, modelIndex, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	}
	return nil, errors.New("image category unavailable")
}

// doesnt support image to image generation
func processLowImageGenerationInput(_ context.Context, r *http.Request, modelIndex int, _ *config.ApiConfig) (replicate.PredictionInput, error) {
	var lowImageGenerationParams rep.LowImageGenerationParams
	decoder := json.NewDecoder(r.Body)
	switch modelIndex {
	case 0:
		lowImageGenerationParams = rep.LowImageGenerationParams{}.NewSdxlLightning4StepImageGenerationInput()
	case 3:
		lowImageGenerationParams = rep.LowImageGenerationParams{}.NewDreamshaperXLTurboImageGenerationInput()
	}

	params := lowImageGenerationParams
	err := decoder.Decode(&params)
	if err != nil {
		return nil, err
	}
	input := replicate.PredictionInput{
		"prompt":              params.Prompt,
		"negative_prompt":     params.NegativePrompt,
		"width":               params.Width,
		"height":              params.Height,
		"scheduler":           params.Scheduler,
		"num_outputs":         params.NumOutputs,
		"guidance_scale":      params.GuidanceScale,
		"num_inference_steps": params.NumInferenceSteps,
	}

	return input, nil
}

// support imagetoimage generation
func processHighImageGenerationInput(ctx context.Context, r *http.Request, modelIndex int, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	var HighImageGenerationParams rep.HighImageGenerationParams
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}

	// Extract the JSON metadata from the form
	switch modelIndex {
	case 1:
		HighImageGenerationParams = rep.HighImageGenerationParams{}.NewRealvisxlV20()
	case 2:
		HighImageGenerationParams = rep.HighImageGenerationParams{}.Newplaygroundv251024()
	case 4:
		HighImageGenerationParams = rep.HighImageGenerationParams{}.NewAstra()
	}

	// Collect form values
	prompt := r.FormValue("prompt")
	if prompt == "" {
		return nil, errors.New("prompt cant be empty")
	}
	HighImageGenerationParams.Prompt = prompt
	//extract and replace default
	utils.SetStringValue(r.FormValue("scheduler"), &HighImageGenerationParams.Scheduler)
	utils.SetStringValue(r.FormValue("negative_prompt"), &HighImageGenerationParams.NegativePrompt)
	utils.SetIntValue(r.FormValue("width"), &HighImageGenerationParams.Width)
	utils.SetIntValue(r.FormValue("height"), &HighImageGenerationParams.Height)
	utils.SetIntValue(r.FormValue("num_outputs"), &HighImageGenerationParams.NumOutputs)
	utils.SetFloatValue(r.FormValue("guidance_scale"), &HighImageGenerationParams.GuidanceScale)
	utils.SetIntValue(r.FormValue("num_inference_steps"), &HighImageGenerationParams.NumInferenceSteps)
	utils.SetFloatValue(r.FormValue("lora_scale"), &HighImageGenerationParams.LoraScale)
	utils.SetBoolValue(r.FormValue("apply_watermark"), &HighImageGenerationParams.ApplyWatermark)
	utils.SetFloatValue(r.FormValue("prompt_strength"), &HighImageGenerationParams.PromptStrength)
	utils.SetIntValue(r.FormValue("seed"), &HighImageGenerationParams.Seed)

	// Handle image file
	imageFile, imageFileHeader, err := r.FormFile("image")
	if err == nil {
		repFile, err := rep.RequestFileToReplicateFile(ctx, imageFileHeader, cfg.ReplicateAPIKey)
		if err != nil {
			return nil, err
		}
		HighImageGenerationParams.ImageFile = repFile
	}
	if imageFile != nil {
		defer imageFile.Close()
	}

	// Handle mask file
	maskFile, maskFileHeader, err := r.FormFile("mask")
	if err == nil {
		repFile, err := rep.RequestFileToReplicateFile(ctx, maskFileHeader, cfg.ReplicateAPIKey)
		if err != nil {
			return nil, err
		}
		HighImageGenerationParams.MaskFile = repFile
	}
	if maskFile != nil {
		defer maskFile.Close()
	}

	input := replicate.PredictionInput{
		"prompt":              HighImageGenerationParams.Prompt,
		"width":               HighImageGenerationParams.Width,
		"height":              HighImageGenerationParams.Height,
		"scheduler":           HighImageGenerationParams.Scheduler,
		"lora_scale":          HighImageGenerationParams.LoraScale,
		"num_outputs":         HighImageGenerationParams.NumOutputs,
		"guidance_scale":      HighImageGenerationParams.GuidanceScale,
		"apply_watermark":     HighImageGenerationParams.ApplyWatermark,
		"negative_prompt":     HighImageGenerationParams.NegativePrompt,
		"prompt_strength":     HighImageGenerationParams.PromptStrength,
		"num_inference_steps": HighImageGenerationParams.NumInferenceSteps,
	}
	if HighImageGenerationParams.ImageFile != nil {
		input["image"] = HighImageGenerationParams.ImageFile
	}
	if HighImageGenerationParams.MaskFile != nil {
		input["mask"] = HighImageGenerationParams.MaskFile
	}
	if HighImageGenerationParams.Seed != nil {
		input["seed"] = *HighImageGenerationParams.Seed
	}

	return input, nil
}
