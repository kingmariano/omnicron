package utils

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv(path string) (string, string, string, string, string, error) {
	err := godotenv.Load(path)
	if err != nil {
		log.Print(err)
	}
	apiKey := os.Getenv("MY_API_KEY")
	if apiKey == "" {
		return "", "", "", "", "", errors.New("unable to get API key")
	}

	grokApiKey := os.Getenv("GROK_API_KEY")
	if grokApiKey == "" {
		return apiKey, "", "", "", "", errors.New("unable to get Grok API key")
	}
	replicateApiKey := os.Getenv("REPLICATE_API_TOKEN")
	if replicateApiKey == "" {
		return apiKey, grokApiKey, "", "", "", errors.New("unable to get Replicate API key")
	}
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		return apiKey, grokApiKey, replicateApiKey, "", "", errors.New("unable to get cloudinary URL")
	}
	port := os.Getenv("PORT")
	if port == "" {
		return apiKey, grokApiKey, replicateApiKey, cloudinaryURL, "", errors.New("unable to get port")
	}
	return apiKey, grokApiKey, replicateApiKey, cloudinaryURL, port, nil

}
