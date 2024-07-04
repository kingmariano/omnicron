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
	"bytes"
	"context"
	replicate "github.com/replicate/replicate-go"
	"io"
	"log"
	"mime/multipart"
)

type CreateReplicatePrediction func(ctx context.Context, env, version string, predictionInput replicate.PredictionInput, webhook *replicate.Webhook, stream bool) (*replicate.Prediction, error)

type Model map[ReplicateModel]CreateReplicatePrediction

var ImageGenModels Model
var ImageUpscaleGenModels Model
var VideoGenModels Model
var TTSGenModels Model
var STTGenModels Model
var MusicGenModels Model

type ReplicateModel struct {
	Name     string
	Version  string
	Category string
}

// init initializes the global variables ImageGenModels, ImageUpscaleGenModels, VideoGenModels,
// TTSGenModels, STTGenModels, and MusicGenModels. It maps the CreatePrediction function to each
// model in the respective categories.
func init() {
	ImageGenModels = make(Model)
	ImageUpscaleGenModels = make(Model)
	VideoGenModels = make(Model)
	TTSGenModels = make(Model)
	STTGenModels = make(Model)
	MusicGenModels = make(Model)
	for _, imagemodels := range ImageModels {
		ImageGenModels[imagemodels] = CreatePrediction
	}
	for _, imageupscalemodels := range ImageUpscaleModels {
		ImageUpscaleGenModels[imageupscalemodels] = CreatePrediction
	}
	for _, videomodels := range VideoModels {
		VideoGenModels[videomodels] = CreatePrediction
	}

	for _, ttsmodels := range TTSModels {
		TTSGenModels[ttsmodels] = CreatePrediction
	}
	for _, sttmodels := range STTModels {
		STTGenModels[sttmodels] = CreatePrediction
	}
	for _, musicmodels := range MusicModels {
		MusicGenModels[musicmodels] = CreatePrediction
	}
}

func NewReplicateClient(token string) (*replicate.Client, error) {
	r8, err := replicate.NewClient(replicate.WithToken(token))
	if err != nil {
		return nil, err
	}
	return r8, nil
}

// CreatePrediction creates a new prediction using the provided token, version, predictionInput, webhook, and stream.
// It uses the NewReplicateClient function to create a new replicate client and then calls the CreatePrediction method on the client.
// After creating the prediction, it waits for the prediction to complete using the Wait method on the client.
// If any error occurs during the process, it returns the error. Otherwise, it logs a success message and returns the prediction.
func CreatePrediction(ctx context.Context, token, version string, predictionInput replicate.PredictionInput, webhook *replicate.Webhook, stream bool) (*replicate.Prediction, error) {
	r8, err := NewReplicateClient(token)
	if err != nil {
		return nil, err
	}

	prediction, err := r8.CreatePrediction(ctx, version, predictionInput, webhook, stream)
	if err != nil {
		return nil, err
	}
	err = r8.Wait(ctx, prediction)
	if err != nil {
		return nil, err
	}
	log.Println("successfully executed prediction")
	return prediction, nil
}

// converts the request file to replicate file
func RequestFileToReplicateFile(ctx context.Context, fileHeader *multipart.FileHeader, token string) (*replicate.File, error) {
	requestFile, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer requestFile.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, requestFile)
	if err != nil {
		return nil, err
	}

	r8, err := NewReplicateClient(token)
	if err != nil {
		return nil, err
	}
	file, err := r8.CreateFileFromBuffer(ctx, buf, &replicate.CreateFileOptions{})

	if err != nil {
		return nil, err
	}
	return file, err
}
