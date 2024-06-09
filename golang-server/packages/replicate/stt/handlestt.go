package stt

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

func processSTTModelInput(STTModel *rep.ReplicateModel, ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	if STTModel.Category == "Low" {
		replicateInput, err := processLowSTTInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	} else if STTModel.Category == "High" {
		replicateInput, err := processHighSTTInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	}
	return nil, errors.New("tts category unavailable")

}

func processLowSTTInput(ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	log.Println("This is  low STT")
	var LowSTTParams rep.LowSTTParams
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}
	LowSTTParams = rep.LowSTTParams{}.Whisper()
	audioFile, _, err := r.FormFile("audio")
	if err != nil {
		return nil, fmt.Errorf("provide audio file: %v", err)
	}
	audioUrl, err := storage.HandleFileUpload(ctx, audioFile, cfg)
	if err != nil {
		return nil, err
	}
	LowSTTParams.AudioURL = audioUrl

	utils.SetStringValue(r.FormValue("transcription"), &LowSTTParams.Transcription)
	utils.SetStringValue(r.FormValue("initial_prompt"), &LowSTTParams.InitialPrompt)
	utils.SetFloatValue(r.FormValue("temperature"), &LowSTTParams.Temperature)
	utils.SetBoolValue(r.FormValue("translate"), &LowSTTParams.Translate)
	utils.SetBoolValue(r.FormValue("condition_on_previous_text"), &LowSTTParams.ConditionOnPreviousText)

	//set InitialPrompt if present

	input := replicate.PredictionInput{
		"audio":                      LowSTTParams.AudioURL,
		"transcription":              LowSTTParams.Transcription,
		"translate":                  LowSTTParams.Translate,
		"temperature":                LowSTTParams.Temperature,
		"condition_on_previous_text": LowSTTParams.ConditionOnPreviousText,
	}

	if LowSTTParams.InitialPrompt != nil {
		input["initial_prompt"] = *LowSTTParams.InitialPrompt
	}
	return input, nil

}
func processHighSTTInput(ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	log.Println("This is  high STT")
	var HighSTTParams rep.HighSTTParams
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}
	HighSTTParams = rep.HighSTTParams{}.InsanelyFastWhisperWithVideo()
	audioFile, _, err := r.FormFile("audio")
	if err != nil {
		return nil, fmt.Errorf("provide audio file: %v", err)
	}
	audioUrl, err := storage.HandleFileUpload(ctx, audioFile, cfg)
	if err != nil {
		return nil, err
	}
	HighSTTParams.AudioURL = audioUrl
	utils.SetStringValue(r.FormValue("task"), &HighSTTParams.Task)
	utils.SetIntValue(r.FormValue("batch_size"), &HighSTTParams.BatchSize)
	utils.SetStringValue(r.FormValue("timestamp"), &HighSTTParams.Timestamp)
	input := replicate.PredictionInput{
		"audio":      HighSTTParams.AudioURL,
		"task":       HighSTTParams.Task,
		"batch_size": HighSTTParams.BatchSize,
		"timestamp":  HighSTTParams.Timestamp,
	}
	return input, nil
}
