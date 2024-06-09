package replicate

import (
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
)

// doesnt support file upload
type LowImageGenerationParams struct {
	Prompt            string   `json:"prompt"`
	Width             *int     `json:"width,omitempty"`
	Height            *int     `json:"height,omitempty"`
	Scheduler         *string  `json:"scheduler,omitempty"`
	NumOutputs        *int     `json:"num_outputs,omitempty"`
	GuidanceScale     *float64 `json:"guidance_scale,omitempty"`
	NegativePrompt    *string  `json:"negative_prompt,omitempty"`
	NumInferenceSteps *int     `json:"num_inference_steps,omitempty"`
}

// support file upload
type HighImageGenerationParams struct {
	Prompt            string   `json:"prompt"`
	Width             *int     `json:"width,omitempty"`
	Height            *int     `json:"height,omitempty"`
	Scheduler         *string  `json:"scheduler,omitempty"`
	NumOutputs        *int     `json:"num_outputs,omitempty"`
	GuidanceScale     *float64 `json:"guidance_scale,omitempty"`
	NegativePrompt    *string  `json:"negative_prompt,omitempty"`
	NumInferenceSteps *int     `json:"num_inference_steps,omitempty"`
	LoraScale         *float64 `json:"lora_scale,omitempty"`
	ImageURL          *string  `json:"image_url,omitempty"`
	MaskURL           *string  `json:"mask_url,omitempty"`
	PromptStrength    *float64 `json:"prompt_strength,omitempty"`
	ApplyWatermark    *bool    `json:"apply_watermark,omitempty"`
	Seed              *int     `json:"seed,omitempty"`
}

type HighImageUpscaleGenerationParams struct {
	ImageURL              string   `json:"image_url"`
	Prompt                *string  `json:"prompt,omitempty"`
	NegativePrompt        *string  `json:"negative_prompt,omitempty"`
	ScaleFactor           *float64 `json:"scale_factor,omitempty"`
	Dynamic               *float64 `json:"dynamic,omitempty"`
	Creativity            *float64 `json:"creativity,omitempty"`
	Resemblance           *float64 `json:"resemblance,omitempty"`
	TilingWidth           *int     `json:"tiling_width,omitempty"`
	TilingHeight          *int     `json:"tiling_height,omitempty"`
	SdModel               *string  `json:"sd_model,omitempty"`
	Scheduler             *string  `json:"scheduler,omitempty"`
	NumInferenceSteps     *int     `json:"num_inference_steps,omitempty"`
	Seed                  *int     `json:"seed,omitempty"`
	Downscaling           *bool    `json:"downscaling,omitempty"`
	DownscalingResolution *int     `json:"downscaling_resolution,omitempty"`
	Sharpen               *float64 `json:"sharpen,omitempty"`
	OutputFormat          *string  `json:"output_format,omitempty"`
}

type LowImageUpscaleGenerationParams struct {
	ImageURL    string   `json:"image_url"`
	Scale       *float64 `json:"scale,omitempty"`
	FaceEnhance *bool    `json:"face_enhance,omitempty"`
}

type HighVideoGenerationParams struct {
	Prompt            string   `json:"prompt"`
	NegativePrompt    *string  `json:"negative_prompt,omitempty"`
	InitVideoURL      *string  `json:"init_video_url,omitempty"`
	InitWeight        *float64 `json:"init_weight,omitempty"`
	NumFrames         *int     `json:"num_frames,omitempty"`
	NumInferenceSteps *int     `json:"num_inferences_steps,omitempty"`
	Width             *int     `json:"width,omitempty"`
	Height            *int     `json:"height,omitempty"`
	GuidanceScale     *float64 `json:"guidance_scale,omitempty"`
	FPS               *int     `json:"fps,omitempty"`
	VideoModel        *string  `json:"video_model,omitempty"`
	BatchSize         *int     `json:"batch_size,omitempty"`
	RemoveWatermark   *bool    `json:"remove_watermark,omitempty"`
}

type LowTTSParams struct {
	Text         string  `json:"text"`
	Speaker      string  `json:"speaker"`
	Language     *string `json:"language,omitempty"`
	CleanupVoice *bool   `json:"cleanup_voice,omitempty"`
}
type MediumTTSParams struct {
	SongInputURL              string   `json:"song_input_url,omitempty"`
	RvcModel                  *string  `json:"rvc_model,omitempty"`
	CustomRvcModelDownloadURL *string  `json:"custom_rvc_model_download_url,omitempty"`
	PitchChange               *string  `json:"pitch_change,omitempty"`
	IndexRate                 *float64 `json:"index_rate,omitempty"`
	FilterRaidus              *int     `json:"filter_raidus,omitempty"`
	RmsMixRate                *float64 `json:"rms_mix_rate,omitempty"`
	PitchDetectionAlgorithm   *string  `json:"pitch_detection_algorithm,omitempty"`
	CrepeHopLength            *int     `json:"crepe_hop_length,omitempty"`
	Protect                   *float64 `json:"protect,omitempty"`
	MainVocalsVolumeChange    *float64 `json:"main_vocals_volume_change,omitempty"`
	BackupVocalsVolumeChange  *float64 `json:"backup_vocals_volume_change,omitempty"`
	InstrumentalVolumeChange  *float64 `json:"instrumental_volume_change,omitempty"`
	PitchChangeAll            *float64 `json:"pitch_change_all,omitempty"`
	ReverbSize                *float64 `json:"reverb_size,omitempty"`
	ReverbWetness             *float64 `json:"reverb_wetness,omitempty"`
	ReverbDryness             *float64 `json:"reverb_dryness,omitempty"`
	ReverbDamping             *float64 `json:"reverb_damping,omitempty"`
	OutputFormat              *string  `json:"output_format,omitempty"`
}

type HighTTSParams struct {
	AudioURL string   `json:"audio_url"`
	Text     string   `json:"text"`
	Language *string  `json:"language,omitempty"`
	Speed    *float64 `json:"speed,omitempty"`
}
type LowSTTParams struct {
	AudioURL                string   `json:"audio_url"`
	Transcription           *string  `json:"transcription,omitempty"`
	Temperature             *float64 `json:"temperature,omitempty"`
	Translate               *bool    `json:"translate,omitempty"`
	InitialPrompt           *string  `json:"initial_prompt,omitempty"`
	ConditionOnPreviousText *bool    `json:"condition_on_previous_text,omitempty"`
}
type HighSTTParams struct {
	AudioURL  string  `json:"audio_url"`
	Task      *string `json:"task,omitempty"`
	BatchSize *int    `json:"batch_size,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
}
type LowMusicGenerationParams struct {
	PromptA           string   `json:"prompt_a"`
	Denoising         *float64 `json:"denoising,omitempty"`
	PromptB           *string  `json:"prompt_b,omitempty"`
	Alpha             *float64 `json:"alpha,omitempty"`
	NumInferenceSteps *int     `json:"num_inference_steps,omitempty"`
	SeedImageID       *string  `json:"seed_image_id,omitempty"`
}
type HighMusicGenerationParams struct {
	Prompt                 string   `json:"prompt"`
	ModelVersion           *string  `json:"model_version,omitempty"`
	InputAudioURL          *string  `json:"input_audio_url,omitempty"`
	Duration               *int     `json:"duration,omitempty"`
	Continuation           *bool    `json:"continuation,omitempty"`
	ContinuationStart      *int     `json:"continuation_start,omitempty"`
	ContinuationEnd        *int     `json:"continuation_end,omitempty"`
	MultiBandDiffusion     *bool    `json:"multi_band_diffusion,omitempty"`
	NormalizationStrategy  *string  `json:"normalization_strategy,omitempty"`
	TopK                   *int     `json:"top_k,omitempty"`
	TopP                   *float64 `json:"top_p,omitempty"`
	Temperature            *float64 `json:"temperature,omitempty"`
	ClassifierFreeGuidance *int     `json:"classifier_free_guidance,omitempty"`
	OutputFormat           *string  `json:"output_format,omitempty"`
}

func (m LowMusicGenerationParams) Riffusion() LowMusicGenerationParams {
	return LowMusicGenerationParams{
		Denoising:         utils.Ptr(0.75),
		Alpha:             utils.Ptr(0.5),
		NumInferenceSteps: utils.Ptr(50),
		SeedImageID:       utils.Ptr("vibes"),
	}
}

func (m HighMusicGenerationParams) MusicGen() HighMusicGenerationParams {
	return HighMusicGenerationParams{
		ModelVersion:           utils.Ptr("stereo-melody-large"),
		Duration:               utils.Ptr(10),
		Continuation:           utils.Ptr(false),
		ContinuationStart:      utils.Ptr(0),
		ContinuationEnd:        utils.Ptr(0),
		MultiBandDiffusion:     utils.Ptr(false),
		NormalizationStrategy:  utils.Ptr("peak"),
		TopK:                   utils.Ptr(50),
		TopP:                   utils.Ptr(0.0),
		Temperature:            utils.Ptr(1.0),
		ClassifierFreeGuidance: utils.Ptr(3),
		OutputFormat:           utils.Ptr("mp3"),
	}
}

func (m LowSTTParams) Whisper() LowSTTParams {
	return LowSTTParams{
		Transcription:           utils.Ptr("plain text"),
		Temperature:             utils.Ptr(0.0),
		Translate:               utils.Ptr(false),
		ConditionOnPreviousText: utils.Ptr(true),
	}
}
func (m HighSTTParams) InsanelyFastWhisperWithVideo() HighSTTParams {
	return HighSTTParams{
		Task:      utils.Ptr("transcribe"),
		BatchSize: utils.Ptr(64),
		Timestamp: utils.Ptr("chunk"),
	}

}

func (m LowTTSParams) XTTSV2() LowTTSParams {
	return LowTTSParams{
		Language:     utils.Ptr("en"),
		CleanupVoice: utils.Ptr(false),
	}
}

func (m MediumTTSParams) RealisticVoiceCloning() MediumTTSParams {
	return MediumTTSParams{
		RvcModel:                 utils.Ptr("Squidward"),
		PitchChange:              utils.Ptr("no-change"),
		IndexRate:                utils.Ptr(0.5),
		FilterRaidus:             utils.Ptr(3),
		RmsMixRate:               utils.Ptr(0.25),
		PitchDetectionAlgorithm:  utils.Ptr("rmvpe"),
		CrepeHopLength:           utils.Ptr(128),
		Protect:                  utils.Ptr(0.33),
		MainVocalsVolumeChange:   utils.Ptr(10.1),
		BackupVocalsVolumeChange: utils.Ptr(0.0),
		InstrumentalVolumeChange: utils.Ptr(0.0),
		PitchChangeAll:           utils.Ptr(0.0),
		ReverbSize:               utils.Ptr(0.15),
		ReverbWetness:            utils.Ptr(0.2),
		ReverbDryness:            utils.Ptr(0.8),
		ReverbDamping:            utils.Ptr(0.7),
		OutputFormat:             utils.Ptr("mp3"),
	}
}
func (m HighTTSParams) OpenVoice() HighTTSParams {
	return HighTTSParams{
		Language: utils.Ptr("EN_NEWEST"),
		Speed:    utils.Ptr(1.0),
	}
}

func (m HighVideoGenerationParams) ZeroscopeV2Xl() HighVideoGenerationParams {
	return HighVideoGenerationParams{
		NegativePrompt:    utils.Ptr("blurred, noisy, washed out, distorted, broken, overly dark, low resolution, excessive blue tones, overexposed, unnatural colors, overly saturated, cluttered, pixelated, abstract"),
		InitWeight:        utils.Ptr(0.5),
		NumFrames:         utils.Ptr(24),
		NumInferenceSteps: utils.Ptr(50),
		Width:             utils.Ptr(1024),
		Height:            utils.Ptr(1024),
		GuidanceScale:     utils.Ptr(17.5),
		FPS:               utils.Ptr(10),
		VideoModel:        utils.Ptr("xl"),
		BatchSize:         utils.Ptr(1),
		RemoveWatermark:   utils.Ptr(false),
	}
}

func (m LowImageUpscaleGenerationParams) RealEsrgan() LowImageUpscaleGenerationParams {
	return LowImageUpscaleGenerationParams{
		Scale:       utils.Ptr(4.0),
		FaceEnhance: utils.Ptr(false),
	}
}

func (m HighImageUpscaleGenerationParams) NewClarityUpscaler() HighImageUpscaleGenerationParams {
	return HighImageUpscaleGenerationParams{
		Prompt:                utils.Ptr("masterpiece, best quality, highres, <lora:more_details:0.5> <lora:SDXLrender_v2.0:1>"),
		NegativePrompt:        utils.Ptr("(worst quality, low quality, normal quality:2) JuggernautNegative-neg"),
		ScaleFactor:           utils.Ptr(2.0),
		Dynamic:               utils.Ptr(6.0),
		Creativity:            utils.Ptr(0.35),
		Resemblance:           utils.Ptr(0.6),
		TilingWidth:           utils.Ptr(112),
		TilingHeight:          utils.Ptr(144),
		SdModel:               utils.Ptr("juggernaut_reborn.safetensors [338b85bc4f]"),
		Scheduler:             utils.Ptr("DPM++ 3M SDE Karras"),
		NumInferenceSteps:     utils.Ptr(18),
		Seed:                  utils.Ptr(1337),
		Downscaling:           utils.Ptr(false),
		DownscalingResolution: utils.Ptr(768),
		Sharpen:               utils.Ptr(0.0),
		OutputFormat:          utils.Ptr("png"),
	}
}

func (m LowImageGenerationParams) NewSdxlLightning4StepImageGenerationInput() LowImageGenerationParams {
	return LowImageGenerationParams{
		Width:             utils.Ptr(1024),
		Height:            utils.Ptr(1024),
		Scheduler:         utils.Ptr("K_EULER"),
		NumOutputs:        utils.Ptr(1),
		GuidanceScale:     utils.Ptr(0.0),
		NegativePrompt:    utils.Ptr("worst quality, low quality, bad anatomy, incorrect perspective"),
		NumInferenceSteps: utils.Ptr(4),
	}

}

func (m LowImageGenerationParams) NewDreamshaperXLTurboImageGenerationInput() LowImageGenerationParams {
	return LowImageGenerationParams{
		Width:             utils.Ptr(1024),
		Height:            utils.Ptr(1024),
		Scheduler:         utils.Ptr("K_EULER"),
		NumOutputs:        utils.Ptr(1),
		GuidanceScale:     utils.Ptr(2.0),
		NegativePrompt:    utils.Ptr("ugly, deformed, noisy, blurry, low contrast, text, watermark, logo, low resolution, bad anatomy, bad proportions, bad lighting, overexposed, underexposed, jpeg artifacts, pixelated, out of focus, cartoon, 3d render, unrealistic, fake, distorted, unnatural, poorly drawn, incorrect perspective, disfigured, messy, cluttered, low detail, poorly rendered, over saturated, washed out"),
		NumInferenceSteps: utils.Ptr(6),
	}
}

func (m HighImageGenerationParams) NewRealvisxlV20() HighImageGenerationParams {
	return HighImageGenerationParams{
		NegativePrompt:    utils.Ptr("ugly, deformed, noisy, blurry, low contrast, text, watermark, logo, low resolution, bad anatomy, bad proportions, bad lighting, overexposed, underexposed, jpeg artifacts, pixelated, out of focus, cartoon, 3d render, unrealistic, fake, distorted, unnatural, poorly drawn, incorrect perspective, disfigured, messy, cluttered, low detail, poorly rendered, over saturated, washed out"),
		Width:             utils.Ptr(1024),
		Height:            utils.Ptr(1024),
		NumOutputs:        utils.Ptr(1),
		GuidanceScale:     utils.Ptr(7.0),
		NumInferenceSteps: utils.Ptr(40),
		PromptStrength:    utils.Ptr(0.8),
		Scheduler:         utils.Ptr("DPMSolverMultistep"),
		ApplyWatermark:    utils.Ptr(false),
		LoraScale:         utils.Ptr(0.6),
	}
}

func (m HighImageGenerationParams) Newplaygroundv251024() HighImageGenerationParams {
	return HighImageGenerationParams{
		NegativePrompt:    utils.Ptr("ugly, deformed, noisy, blurry, low contrast, text, watermark, logo, low resolution, bad anatomy, bad proportions, bad lighting, overexposed, underexposed, jpeg artifacts, pixelated, out of focus, cartoon, 3d render, unrealistic, fake, distorted, unnatural, poorly drawn, incorrect perspective, disfigured, messy, cluttered, low detail, poorly rendered, over saturated, washed out"),
		Width:             utils.Ptr(1024),
		Height:            utils.Ptr(1024),
		NumOutputs:        utils.Ptr(1),
		GuidanceScale:     utils.Ptr(3.0),
		NumInferenceSteps: utils.Ptr(25),
		PromptStrength:    utils.Ptr(0.8),
		Scheduler:         utils.Ptr("DPMSolver++"),
		ApplyWatermark:    utils.Ptr(false),
		LoraScale:         utils.Ptr(0.6),
	}
}

func (m HighImageGenerationParams) NewAstra() HighImageGenerationParams {
	return HighImageGenerationParams{
		NegativePrompt:    utils.Ptr("deformed iris, deformed pupils, semi-realistic, text, cropped, out of frame, worst quality, low quality, jpeg artifacts, ugly, duplicate, morbid, mutilated, extra fingers, mutated hands, poorly drawn hands, poorly drawn face, mutation, deformed, blurry, dehydrated, bad anatomy, bad proportions, extra limbs, cloned face, disfigured, gross proportions, malformed limbs, missing arms, missing legs, extra arms, extra legs, fused fingers, too many fingers, long neck, blurry, low quality , bad quality , Not detailed, watermark, deformed figures, lack of details, bad anatomy, blurry, extra arms, extra fingers, poorly drawn hands, disfigured, tiling, deformed, mutated ,ugly, disfigured, low quality, blurry ,distorted, blur, smooth, low-quality, warm, haze, over-saturated, high-contrast, out of focus, dark, worst quality, low quality"),
		Width:             utils.Ptr(1024),
		Height:            utils.Ptr(1024),
		NumOutputs:        utils.Ptr(1),
		GuidanceScale:     utils.Ptr(7.5),
		NumInferenceSteps: utils.Ptr(50),
		PromptStrength:    utils.Ptr(0.8),
		Scheduler:         utils.Ptr("K_EULER"),
		ApplyWatermark:    utils.Ptr(false),
		LoraScale:         utils.Ptr(0.6),
	}
}
