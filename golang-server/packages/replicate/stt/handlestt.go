package stt

import (
	"context"
	"errors"
	"fmt"
	"github.com/kingmariano/omnicron-backendsever/golang-server/config"
	rep "github.com/kingmariano/omnicron-backendsever/golang-server/packages/replicate"
	"github.com/kingmariano/omnicron-backendsever/golang-server/utils"
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
	_, audioFileHeader, err := r.FormFile("audio")
	if err != nil {
		return nil, fmt.Errorf("provide audio file: %v", err)
	}
	repFile, err := rep.RequestFileToReplicateFile(ctx, audioFileHeader, cfg.ReplicateAPIKey)
	if err != nil {
		return nil, err
	}
	LowSTTParams.AudioFile = repFile

	utils.SetStringValue(r.FormValue("transcription"), &LowSTTParams.Transcription)
	utils.SetStringValue(r.FormValue("initial_prompt"), &LowSTTParams.InitialPrompt)
	utils.SetFloatValue(r.FormValue("temperature"), &LowSTTParams.Temperature)
	utils.SetBoolValue(r.FormValue("translate"), &LowSTTParams.Translate)
	utils.SetBoolValue(r.FormValue("condition_on_previous_text"), &LowSTTParams.ConditionOnPreviousText)

	//set InitialPrompt if present

	input := replicate.PredictionInput{
		"audio":                      LowSTTParams.AudioFile,
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
	audioFile, audioFileHeader, err := r.FormFile("audio")
	if err == nil {
		repFile, err := rep.RequestFileToReplicateFile(ctx, audioFileHeader, cfg.ReplicateAPIKey)
		if err != nil {
			return nil, err
		}
		HighSTTParams.AudioFile = repFile
	}
	if audioFile != nil {
		defer audioFile.Close()
	}
	utils.SetStringValue(r.FormValue("url"), &HighSTTParams.URL)
	utils.SetStringValue(r.FormValue("task"), &HighSTTParams.Task)
	utils.SetIntValue(r.FormValue("batch_size"), &HighSTTParams.BatchSize)
	utils.SetStringValue(r.FormValue("timestamp"), &HighSTTParams.Timestamp)
	input := replicate.PredictionInput{
		"task":       HighSTTParams.Task,
		"batch_size": HighSTTParams.BatchSize,
		"timestamp":  HighSTTParams.Timestamp,
	}
	if HighSTTParams.AudioFile != nil && HighSTTParams.URL != nil {
		return nil, errors.New("audio file and url can't be present at the same time")
	} else if HighSTTParams.AudioFile == nil && HighSTTParams.URL == nil {
		return nil, errors.New("either audio file or url must be present at the")
	} else if HighSTTParams.AudioFile != nil {
		input["audio"] = HighSTTParams.AudioFile
	} else if HighSTTParams.URL != nil {
		input["url"] = *HighSTTParams.URL
	}
	return input, nil
}
