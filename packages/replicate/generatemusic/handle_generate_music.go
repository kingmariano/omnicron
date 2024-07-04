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
package generatemusic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/kingmariano/omnicron/config"
	rep "github.com/kingmariano/omnicron/packages/replicate"
	"github.com/kingmariano/omnicron/utils"
	replicate "github.com/replicate/replicate-go"
	"log"
	"net/http"
)

func processMusicModelInput(MusicModel *rep.ReplicateModel, ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	if MusicModel.Category == "Low" {
		replicateInput, err := processLowMusicGenInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	} else if MusicModel.Category == "High" {
		replicateInput, err := processHighMusicGenInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	}
	return nil, errors.New("music model category unavailable")
}

func processLowMusicGenInput(_ context.Context, r *http.Request, _ *config.ApiConfig) (replicate.PredictionInput, error) {
	log.Println("This is low music model generation")
	var LowMusicGenerationModelsParams rep.LowMusicGenerationParams
	decoder := json.NewDecoder(r.Body)
	LowMusicGenerationModelsParams = rep.LowMusicGenerationParams{}.Riffusion()

	err := decoder.Decode(&LowMusicGenerationModelsParams)
	if err != nil {
		return nil, err
	}

	input := replicate.PredictionInput{
		"prompt_a":            LowMusicGenerationModelsParams.PromptA,
		"denoising":           LowMusicGenerationModelsParams.Denoising,
		"alpha":               LowMusicGenerationModelsParams.Alpha,
		"num_inference_steps": LowMusicGenerationModelsParams.NumInferenceSteps,
		"seed_image_id":       LowMusicGenerationModelsParams.SeedImageID,
	}
	if LowMusicGenerationModelsParams.PromptB != nil {
		input["prompt_b"] = *LowMusicGenerationModelsParams.PromptB
	}
	return input, nil

}
func processHighMusicGenInput(ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	log.Println("This is high music model generation")
	var HighMusicGenerationParams rep.HighMusicGenerationParams
	prompt := r.FormValue("prompt")
	if prompt == "" {
		return nil, errors.New("prompt field parameter is required")
	}
	HighMusicGenerationParams = rep.HighMusicGenerationParams{}.MusicGen()
	HighMusicGenerationParams.Prompt = prompt
	utils.SetStringValue(r.FormValue("model_version"), &HighMusicGenerationParams.ModelVersion)
	utils.SetIntValue(r.FormValue("duration"), &HighMusicGenerationParams.Duration)
	utils.SetBoolValue(r.FormValue("continuation"), &HighMusicGenerationParams.Continuation)
	utils.SetIntValue(r.FormValue("continuation_start"), &HighMusicGenerationParams.ContinuationStart)
	utils.SetIntValue(r.FormValue("continuation_end"), &HighMusicGenerationParams.ContinuationEnd)
	utils.SetBoolValue(r.FormValue("multi_band_diffusion"), &HighMusicGenerationParams.MultiBandDiffusion)
	utils.SetStringValue(r.FormValue("normalization_strategy"), &HighMusicGenerationParams.NormalizationStrategy)
	utils.SetIntValue(r.FormValue("top_k"), &HighMusicGenerationParams.TopK)
	utils.SetFloatValue(r.FormValue("top_p"), &HighMusicGenerationParams.TopP)
	utils.SetFloatValue(r.FormValue("temperature"), &HighMusicGenerationParams.Temperature)
	utils.SetIntValue(r.FormValue("classifier_free_guidance"), &HighMusicGenerationParams.ClassifierFreeGuidance)
	utils.SetStringValue(r.FormValue("output_format"), &HighMusicGenerationParams.OutputFormat)
	inputAudioFile, inputAudioFileHeader, err := r.FormFile("input_audio")
	if err == nil {
		repFile, err := rep.RequestFileToReplicateFile(ctx, inputAudioFileHeader, cfg.ReplicateAPIKey)
		if err != nil {
			return nil, err
		}
		HighMusicGenerationParams.InputAudioFile = repFile
	}
	if inputAudioFile != nil {
		defer inputAudioFile.Close()
	}
	input := replicate.PredictionInput{
		"prompt":                   HighMusicGenerationParams.Prompt,
		"model_version":            HighMusicGenerationParams.ModelVersion,
		"duration":                 HighMusicGenerationParams.Duration,
		"continuation":             HighMusicGenerationParams.Continuation,
		"continuation_start":       HighMusicGenerationParams.ContinuationStart,
		"continuation_end":         HighMusicGenerationParams.ContinuationEnd,
		"multi_band_diffusion":     HighMusicGenerationParams.MultiBandDiffusion,
		"normalization_strategy":   HighMusicGenerationParams.NormalizationStrategy,
		"top_k":                    HighMusicGenerationParams.TopK,
		"top_p":                    HighMusicGenerationParams.TopP,
		"temperature":              HighMusicGenerationParams.Temperature,
		"classifier_free_guidance": HighMusicGenerationParams.ClassifierFreeGuidance,
		"output_format":            HighMusicGenerationParams.OutputFormat,
	}
	if HighMusicGenerationParams.InputAudioFile != nil {
		input["input_audio"] = HighMusicGenerationParams.InputAudioFile
	}
	return input, nil

}
