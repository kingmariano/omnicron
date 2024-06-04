package storage

import (
	"context"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	cldconfig "github.com/cloudinary/cloudinary-go/config"
	"go.uber.org/zap"
)

// upload file data and get public direct url to the file
func HandleFileUpload(ctx context.Context, file interface{}, cfg *config.ApiConfig) (string, error) {
	var URLString string
	cld, err := cloudinary.NewFromURL(cfg.CloudinaryUrl)
	if err != nil {
		return "", err
	}
	var zapLogger, _ = zap.NewDevelopment()
	cld.Logger.Writer = zapLogger.Sugar().With("source", "cloudinary")
	cloudinaryConfig, err := cldconfig.NewFromURL(cfg.CloudinaryUrl)
	if err != nil {
		return "", err
	}
	upload, err := uploader.NewWithConfiguration(cloudinaryConfig)
	if err != nil {
		return "", err
	}
	uploadResult, err := upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return "", err
	}
	URLString = uploadResult.URL
	return URLString, nil

}
