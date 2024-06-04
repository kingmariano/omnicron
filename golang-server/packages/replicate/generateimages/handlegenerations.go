package generateimages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	rep "github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate"
	"github.com/charlesozo/omnicron-backendsever/golang-server/storage"
	replicate "github.com/replicate/replicate-go"
	"log"
	"net/http"
	"strconv"
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
	log.Print("my prompt is ", prompt)
	setStringValue(r.FormValue("scheduler"), &HighImageGenerationParams.Scheduler)
	setStringValue(r.FormValue("negative_prompt"), &HighImageGenerationParams.NegativePrompt)
	setIntValue(r.FormValue("width"), &HighImageGenerationParams.Width)
	setIntValue(r.FormValue("height"), &HighImageGenerationParams.Height)
	setIntValue(r.FormValue("num_outputs"), &HighImageGenerationParams.NumOutputs)
	setFloatValue(r.FormValue("guidance_scale"), &HighImageGenerationParams.GuidanceScale)
	setIntValue(r.FormValue("num_inference_steps"), &HighImageGenerationParams.NumInferenceSteps)
	setFloatValue(r.FormValue("lora_scale"), &HighImageGenerationParams.LoraScale)
	setBoolValue(r.FormValue("apply_watermark"), &HighImageGenerationParams.ApplyWatermark)
	setFloatValue(r.FormValue("prompt_strength"), &HighImageGenerationParams.PromptStrength)
	setIntValue(r.FormValue("seed"), &HighImageGenerationParams.Seed)
	if prompt == "" {
		return nil, errors.New("prompt cant be empty")
	}
	HighImageGenerationParams.Prompt = prompt
	// Handle image file
	imageFile, _, err := r.FormFile("image")
	if err == nil {
		imageUrl, err := storage.HandleFileUpload(ctx, imageFile, cfg)
		if err != nil {
			return nil, err
		}
		HighImageGenerationParams.ImageURL = &imageUrl
	}
	if imageFile != nil {
		defer imageFile.Close()
	}

	// Handle mask file
	maskFile, _, err := r.FormFile("mask")
	if err == nil {
		maskUrl, err := storage.HandleFileUpload(ctx, maskFile, cfg)
		if err != nil {
			return nil, err
		}
		HighImageGenerationParams.MaskURL = &maskUrl
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
	if HighImageGenerationParams.ImageURL != nil {
		input["image"] = *HighImageGenerationParams.ImageURL
	}
	if HighImageGenerationParams.MaskURL != nil {
		input["mask"] = *HighImageGenerationParams.MaskURL
	}
	if HighImageGenerationParams.Seed != nil {
		input["seed"] = *HighImageGenerationParams.Seed
	}
	log.Printf("this is input for highImage %s", input)

	return input, nil
}

func setIntValue(value string, param **int) {
	if value != "" {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			*param = &intValue
		}
	}
}

func setFloatValue(value string, param **float64) {
	if value != "" {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			*param = &floatValue
		}
	}
}

func setBoolValue(value string, param **bool) {
	if value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			*param = &boolValue
		}
	}
}

func setStringValue(value string, param **string) {
	if value != "" {
		*param = &value
	}
}
