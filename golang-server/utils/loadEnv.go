package utils

import (
	"log"
	"errors"
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv(path  string) (string, string, string, error) {
	err := godotenv.Load(path)
	if err != nil{
		log.Print(err)
	}
	apiKey := os.Getenv("MY_API_KEY")
	if apiKey == "" {
		return "", "", "", errors.New("unable to get API key")
	}

	grokApiKey := os.Getenv("GROK_API_KEY")
	if grokApiKey == "" {
		return apiKey, "", "", errors.New("unable to get Grok API key")
	}
	port := os.Getenv("PORT")
	if port == "" {
		return apiKey, grokApiKey, "", errors.New("unable to get port")
	}

	return apiKey, grokApiKey, port, nil
}
