package utils

import (
	"context"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	cldconfig "github.com/cloudinary/cloudinary-go/config"
	"go.uber.org/zap"
)

// HandleFileUpload is a function that handles file uploads to Cloudinary.
// It  uploads the file to Cloudinary and returns the URL of the uploaded file.
//
// If any error occurs during the process, the function returns an empty string and the error.
func HandleFileUpload(ctx context.Context, file interface{}, cloudinaryURL string) (string, error) {
	var URLString string
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return "", err
	}
	var zapLogger, _ = zap.NewDevelopment()
	cld.Logger.Writer = zapLogger.Sugar().With("source", "cloudinary")
	cloudinaryConfig, err := cldconfig.NewFromURL(cloudinaryURL)
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
