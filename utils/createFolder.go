package utils

import (
	"fmt"
	"github.com/google/uuid"
	"os"
)

// set the default path to where the video or audio will be downloaded
var BasePath string = "./downloads"
var OutputName string = "youtube"

func CreateUniqueFolder(basePath string) (string, error) {
	uniqueFolder := basePath + uuid.New().String()

	err := os.MkdirAll(uniqueFolder, 0750)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}

	return uniqueFolder, nil
}

func DeleteFolder(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("error deleting directory: %v", err)
	}
	return nil
}
