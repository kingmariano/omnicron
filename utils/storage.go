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
// SOFTWARE.

package utils

import (
	"context"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	cldconfig "github.com/cloudinary/cloudinary-go/config"
	"go.uber.org/zap"
)

// HandleFileUpload is a function that handles file uploads to Cloudinary.
// It  uploads the file to Cloudinary and returns the Direct URL of the uploaded file.
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
