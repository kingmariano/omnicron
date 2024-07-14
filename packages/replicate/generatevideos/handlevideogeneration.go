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
package generatevideos

import (
	"context"
	"errors"
	"fmt"
	"github.com/kingmariano/omnicron/config"
	rep "github.com/kingmariano/omnicron/packages/replicate"
	"github.com/kingmariano/omnicron/utils"
	replicate "github.com/replicate/replicate-go"
	"net/http"
)

func processVideoModelInput(imageModel *rep.ReplicateModel, ctx context.Context, r *http.Request, cfg *config.APIConfig) (replicate.PredictionInput, error) {
	if imageModel.Category == "High" {
		replicateInput, err := processHighVideoInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	}
	return nil, errors.New("image category unavailable")
}

func processHighVideoInput(ctx context.Context, r *http.Request, cfg *config.APIConfig) (replicate.PredictionInput, error) {
	var HighVideoGenerationParams rep.HighVideoGenerationParams
	err := r.ParseMultipartForm(50 << 20) // 50MB
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}

	HighVideoGenerationParams = rep.HighVideoGenerationParams{}.ZeroscopeV2Xl()
	prompt := r.FormValue("prompt")
	if prompt == "" {
		return nil, errors.New("prompt cant be empty")
	}
	HighVideoGenerationParams.Prompt = prompt
	//extract and replace default
	utils.SetStringValue(r.FormValue("negative_prompt"), &HighVideoGenerationParams.NegativePrompt)
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
	videoFile, videoFileHeader, err := r.FormFile("video")
	if err == nil {
		repFile, err := rep.RequestFileToReplicateFile(ctx, videoFileHeader, cfg.ReplicateAPIKey)
		if err != nil {
			return nil, err
		}
		HighVideoGenerationParams.InitVideoFile = repFile
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
	if HighVideoGenerationParams.InitVideoFile != nil {
		input["init_video"] = *HighVideoGenerationParams.InitVideoFile
	}
	return input, nil

}
