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

