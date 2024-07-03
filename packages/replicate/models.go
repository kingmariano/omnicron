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
package replicate

import (
	"errors"
)

var ImageModels = []ReplicateModel{
	{
		Name:     "bytedance/sdxl-lightning-4step",
		Version:  "5f24084160c9089501c1b3545d9be3c27883ae2239b6f412990e82d4a6210f8f",
		Category: "Low",
	},
	{
		Name:     "lucataco/realvisxl-v2.0",
		Version:  "7d6a2f9c4754477b12c14ed2a58f89bb85128edcdd581d24ce58b6926029de08",
		Category: "High",
	},
	{
		Name:     "playgroundai/playground-v2.5-1024px-aesthetic",
		Version:  "a45f82a1382bed5c7aeb861dac7c7d191b0fdf74d8d57c4a0e6ed7d4d0bf7d24",
		Category: "High",
	},
	{
		Name:     "lucataco/dreamshaper-xl-turbo",
		Version:  "0a1710e0187b01a255302738ca0158ff02a22f4638679533e111082f9dd1b615",
		Category: "Low",
	},
	{
		Name:     "lorenzomarines/astra",
		Version:  "6ce68112bcaefc7273692243c933c2dcbe0307757a932fede1ca5e12956e0029",
		Category: "High",
	},
}

var ImageUpscaleModels = []ReplicateModel{
	{
		Name:     "nightmareai/real-esrgan",
		Version:  "350d32041630ffbe63c8352783a26d94126809164e54085352f8326e53999085",
		Category: "Low",
	},
	{
		Name:     "philz1337x/clarity-upscaler",
		Version:  "f11a4727f8f995d2795079196ebda1bcbc641938e032154f46488fc3e760eb79",
		Category: "High",
	},
}
var VideoModels = []ReplicateModel{
	{
		Name:     "anotherjesse/zeroscope-v2-xl",
		Version:  "9f747673945c62801b13b84701c783929c0ee784e4748ec062204894dda1a351",
		Category: "High",
	},
}

var TTSModels = []ReplicateModel{
	{
		Name:     "lucataco/xtts-v2",
		Version:  "684bc3855b37866c0c65add2ff39c78f3dea3f4ff103a436465326e0f438d55e",
		Category: "Low",
	},
	{
		Name:     "zsxkib/realistic-voice-cloning",
		Version:  "0a9c7c558af4c0f20667c1bd1260ce32a2879944a0b9e44e1398660c077b1550",
		Category: "Medium",
	},
	{
		Name:     "chenxwh/openvoice",
		Version:  "d548923c9d7fc9330a3b7c7f9e2f91b2ee90c83311a351dfcd32af353799223d",
		Category: "High",
	},
}

var STTModels = []ReplicateModel{
	{
		Name:     "openai/whisper",
		Version:  "4d50797290df275329f202e48c76360b3f22b08d28c196cbc54600319435f8d2",
		Category: "Low",
	},
	{
		Name:     "turian/insanely-fast-whisper-with-video",
		Version:  "4f41e90243af171da918f04da3e526b2c247065583ea9b757f2071f573965408",
		Category: "High",
	},
}
var MusicModels = []ReplicateModel{
	{
		Name:     "riffusion/riffusion",
		Version:  "8cf61ea6c56afd61d8f5b9ffd14d7c216c0a93844ce2d82ac1c9ecc9c7f24e05",
		Category: "Low",
	},
	{
		Name:     "meta/musicgen",
		Version:  "671ac645ce5e552cc63a54a2bbff63fcf798043055d2dac5fc9e36a837eedcfb",
		Category: "High",
	},
}

func GetModelByName(name string, modelList []ReplicateModel) (*ReplicateModel, error) {
	for _, model := range modelList {
		if model.Name == name {
			return &model, nil
		}
	}
	return nil, errors.New("model not found")
}

func GetModelIndex(modelName string, modelList []ReplicateModel) int {
	for i, model := range modelList {
		if model.Name == modelName {
			return i
		}
	}
	return -1
}
