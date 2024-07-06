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
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// this function loads the environment varaibles from the path
func LoadEnv(path string) (string, string, string, string, string, string, error) {
	err := godotenv.Load(path)
	if err != nil {
		log.Print(err)
	}
	apiKey := os.Getenv("MY_API_KEY")
	if apiKey == "" {
		return "", "", "", "", "", "", errors.New("unable to get API key")
	}

	grokApiKey := os.Getenv("GROK_API_KEY")
	if grokApiKey == "" {
		return apiKey, "", "", "", "", "", errors.New("unable to get Grok API key")
	}
	replicateApiKey := os.Getenv("REPLICATE_API_TOKEN")
	if replicateApiKey == "" {
		return apiKey, grokApiKey, "", "", "", "", errors.New("unable to get Replicate API key")
	}
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		return apiKey, grokApiKey, replicateApiKey, "", "", "", errors.New("unable to get cloudinary URL")
	}
	port := os.Getenv("PORT")
	if port == "" {
		return apiKey, grokApiKey, replicateApiKey, cloudinaryURL, "", "", errors.New("unable to get port")
	}
	youtubeDeveloperKey := os.Getenv("YOUTUBE_DEVELOPER_KEY")
	if youtubeDeveloperKey == "" {
		return apiKey, grokApiKey, replicateApiKey, cloudinaryURL, port, "", errors.New("unable to getyoutube developer")
	}
	log.Printf("all env is set apikey: %s, repicateApiKey: %s, cloudinaryURL: %s", apiKey, replicateApiKey, cloudinaryURL)
	
	return apiKey, grokApiKey, replicateApiKey, cloudinaryURL, port, youtubeDeveloperKey, nil

}
