package tts

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

func processTTSModelInput(TTSModel *rep.ReplicateModel, ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	if TTSModel.Category == "Low" {
		replicateInput, err := processLowTTSInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	} else if TTSModel.Category == "Medium" {
		replicateInput, err := processMediumTTSInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	} else if TTSModel.Category == "High" {
		replicateInput, err := processHighTTSInput(ctx, r, cfg)
		if err != nil {
			return nil, err
		}
		return replicateInput, nil
	}
	return nil, errors.New("tts category unavailable")
}

func processLowTTSInput(ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	log.Println("This is  low TTS")
	var LowTTSParams rep.LowTTSParams
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}
	LowTTSParams = rep.LowTTSParams{}.XTTSV2()
	_, audioFileHeader, err := r.FormFile("audio")
	if err != nil {
		return nil, fmt.Errorf("provide audio file: %v", err)
	}
	repFile, err := rep.RequestFileToReplicateFile(ctx, audioFileHeader, cfg.ReplicateAPIKey)
	if err != nil {
		return nil, err
	}
	text := r.FormValue("text")
	if text == "" {
		return nil, errors.New("please provide the text parameter")
	}
	LowTTSParams.SpeakerFile = repFile
	LowTTSParams.Text = text
	utils.SetStringValue(r.FormValue("language"), &LowTTSParams.Language)
	utils.SetBoolValue(r.FormValue("cleanup_voice"), &LowTTSParams.CleanupVoice)
	input := replicate.PredictionInput{
		"text":          LowTTSParams.Text,
		"speaker":       LowTTSParams.SpeakerFile,
		"language":      LowTTSParams.Language,
		"cleanup_voice": LowTTSParams.CleanupVoice,
	}
	return input, nil
}

func processMediumTTSInput(ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	log.Println("This is  medium TTS")
	var MediumTTSParams rep.MediumTTSParams
	err := r.ParseMultipartForm(50 << 20) // 50MB
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}
	MediumTTSParams = rep.MediumTTSParams{}.RealisticVoiceCloning()
	_, audioFileHeader, err := r.FormFile("audio")
	if err != nil {
		return nil, fmt.Errorf("provide audio file: %v", err)
	}
	repFile, err := rep.RequestFileToReplicateFile(ctx, audioFileHeader, cfg.ReplicateAPIKey)
	if err != nil {
		return nil, err
	}
	MediumTTSParams.SongInputFile = repFile
	utils.SetStringValue(r.FormValue("rvc_model"), &MediumTTSParams.RvcModel)
	utils.SetStringValue(r.FormValue("pitch_change"), &MediumTTSParams.PitchChange)
	utils.SetFloatValue(r.FormValue("index_rate"), &MediumTTSParams.IndexRate)
	utils.SetIntValue(r.FormValue("filter_raidus"), &MediumTTSParams.FilterRaidus)
	utils.SetFloatValue(r.FormValue("rms_mix_rate"), &MediumTTSParams.RmsMixRate)
	utils.SetStringValue(r.FormValue("pitch_detection_algorithm"), &MediumTTSParams.PitchDetectionAlgorithm)
	utils.SetIntValue(r.FormValue("crepe_hop_length"), &MediumTTSParams.CrepeHopLength)
	utils.SetFloatValue(r.FormValue("protect"), &MediumTTSParams.Protect)
	utils.SetFloatValue(r.FormValue("main_vocals_volume_change"), &MediumTTSParams.MainVocalsVolumeChange)
	utils.SetFloatValue(r.FormValue("backup_vocals_volume_change"), &MediumTTSParams.BackupVocalsVolumeChange)
	utils.SetFloatValue(r.FormValue("instrumental_volume_change"), &MediumTTSParams.InstrumentalVolumeChange)
	utils.SetFloatValue(r.FormValue("pitch_change_all"), &MediumTTSParams.PitchChangeAll)
	utils.SetFloatValue(r.FormValue("reverb_size"), &MediumTTSParams.ReverbSize)
	utils.SetFloatValue(r.FormValue("reverb_wetness"), &MediumTTSParams.ReverbWetness)
	utils.SetFloatValue(r.FormValue("reverb_dryness"), &MediumTTSParams.ReverbDryness)
	utils.SetFloatValue(r.FormValue("reverb_damping"), &MediumTTSParams.ReverbDamping)
	utils.SetStringValue(r.FormValue("output_format"), &MediumTTSParams.OutputFormat)

	input := replicate.PredictionInput{
		"protect":                     MediumTTSParams.Protect,
		"rvc_model":                   MediumTTSParams.RvcModel,
		"index_rate":                  MediumTTSParams.IndexRate,
		"song_input":                  MediumTTSParams.SongInputFile,
		"reverb_size":                 MediumTTSParams.ReverbSize,
		"pitch_change":                MediumTTSParams.PitchChange,
		"rms_mix_rate":                MediumTTSParams.RmsMixRate,
		"filter_radius":               MediumTTSParams.FilterRaidus,
		"output_format":               MediumTTSParams.OutputFormat,
		"reverb_damping":              MediumTTSParams.ReverbDamping,
		"reverb_dryness":              MediumTTSParams.ReverbDamping,
		"reverb_wetness":              MediumTTSParams.ReverbWetness,
		"crepe_hop_length":            MediumTTSParams.CrepeHopLength,
		"pitch_change_all":            MediumTTSParams.PitchChangeAll,
		"main_vocals_volume_change":   MediumTTSParams.MainVocalsVolumeChange,
		"pitch_detection_algorithm":   MediumTTSParams.PitchDetectionAlgorithm,
		"instrumental_volume_change":  MediumTTSParams.InstrumentalVolumeChange,
		"backup_vocals_volume_change": MediumTTSParams.BackupVocalsVolumeChange,
	}
	if *MediumTTSParams.RvcModel == "CUSTOM" {
		speechModel := r.FormValue("speech_model")
		if speechModel == "" {
			return nil, errors.New("speech model cant be empty")
		}
		MediumTTSParams.CustomRvcModelDownloadURL = &speechModel

		input["custom_rvc_model_download_url"] = *MediumTTSParams.CustomRvcModelDownloadURL
	}
	return input, nil
}
func processHighTTSInput(ctx context.Context, r *http.Request, cfg *config.ApiConfig) (replicate.PredictionInput, error) {
	log.Println("This is  High TTS")
	var HighTTSParams rep.HighTTSParams
	err := r.ParseMultipartForm(50 << 20) // 50MB
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}
	HighTTSParams = rep.HighTTSParams{}.OpenVoice()
	_, audioFileHeader, err := r.FormFile("audio")
	if err != nil {
		return nil, fmt.Errorf("provide audio file: %v", err)
	}
	repFile, err := rep.RequestFileToReplicateFile(ctx, audioFileHeader, cfg.ReplicateAPIKey)
	if err != nil {
		return nil, err
	}
	text := r.FormValue("text")
	if text == "" {
		return nil, errors.New("please provide the text parameter")
	}
	HighTTSParams.Text = text
	HighTTSParams.AudioFile = repFile
	utils.SetStringValue(r.FormValue("language"), &HighTTSParams.Language)
	utils.SetFloatValue(r.FormValue("Speed"), &HighTTSParams.Speed)
	input := replicate.PredictionInput{
		"text":     HighTTSParams.Text,
		"audio":    HighTTSParams.AudioFile,
		"speed":    HighTTSParams.Speed,
		"language": HighTTSParams.Language,
	}
	return input, nil
}
