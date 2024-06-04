package replicate

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"os"

	replicate "github.com/replicate/replicate-go"
)

type CreateReplicatePrediction func(ctx context.Context, env, version string, predictionInput replicate.PredictionInput, webhook *replicate.Webhook, stream bool) (*replicate.Prediction, error)

type Model map[ReplicateModel]CreateReplicatePrediction

var ImageGenModels Model
var ImageUpscaleGenModels Model

type ReplicateModel struct {
	Name     string
	Version  string
	Category string
}

func init() {
	ImageGenModels = make(Model)
	ImageUpscaleGenModels = make(Model)
	for _, imagemodels := range ImageModels {
		// log.Println(imagemodels)
		ImageGenModels[imagemodels] = CreatePrediction
	}
	for _, imageupscalemodels := range ImageUpscaleModels {
		log.Println(imageupscalemodels)
		ImageUpscaleGenModels[imageupscalemodels] = CreatePrediction
	}

}

func NewReplicateClient(token string) (*replicate.Client, error) {
	r8, err := replicate.NewClient(replicate.WithToken(token))
	if err != nil {
		return nil, err
	}
	return r8, nil
}

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

// convert the request file to replicate file
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

	log.Println("this is file URL", file.URLs["get"])
	if err != nil {
		return nil, err
	}
	jsonData, err := json.MarshalIndent(file, "", " ")
	if err != nil {
		return nil, err
	}

	jsonfile, err := os.Create("file.json")
	if err != nil {

		return nil, err
	}
	defer jsonfile.Close()

	_, err = jsonfile.Write(jsonData)
	if err != nil {
		return nil, err
	}
	log.Println("JSON data written to params.json successfully")

	return file, err
}
