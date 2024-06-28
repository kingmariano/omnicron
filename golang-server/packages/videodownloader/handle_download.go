package videodownloader

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/iawia002/lux/app"
	"github.com/iawia002/lux/downloader"
	"github.com/iawia002/lux/extractors"
)

// extractUrl is a function that extracts video data from a given URL using the lux library.
//
// Parameters:
// url: A string representing the URL of the video to be extracted.
//
// Returns:
// A slice of pointers to extractors.Data, representing the extracted video data.
func extractUrl(url string) ([]*extractors.Data, error) {

	data, err := extractors.Extract(url, extractors.Options{})
	if err != nil {
		return nil, err
	}
	return data, nil
}
func deleteContents(dir string) error {

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())

		// If it's a directory, recursively delete its contents
		if file.IsDir() {
			err = os.RemoveAll(filePath)
		} else {
			err = os.Remove(filePath)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// handleStreamResolution is a function that converts a given video resolution string
// into the corresponding stream identifier used by the lux library.
//
// Parameters:
// resolution: A string representing the desired video resolution.
//
//	Supported values are "1080p", "720p", "480p", "360p", "240p".
//
// Returns:
// A string representing the corresponding stream identifier.
// If the input resolution is not supported, an empty string is returned.
func handleStreamResolution(resolution string) string {
	var stream string
	switch resolution {
	case "1080p":
		stream = "137"
	case "720p":
		stream = "136"
	case "480p":
		stream = "135"
	case "360p":
		stream = "396"
	case "240p":
		stream = "133"
	default:
		stream = ""
	}
	return stream
}

// DownloadVideoData is a function that downloads a video from a given URL,
// with the specified output name, path, and resolution.
// It uses the lux library for extracting video data and downloading the video.
// If an error occurs during the process, it cleans up by removing unnecessary files.
//
// Parameters:
// url: The URL of the video to be downloaded.
// outputName: The name of the output file.
// outputPath: The path where the output file will be saved.
// resolution: The desired resolution of the video.
//
// Returns:
// An error if any error occurs during the process, otherwise nil.
func DownloadVideoData(url string, outputName string, outputPath string, resolution string) (string, error) {
	data, err := extractUrl(url)
	if err != nil {
		return "", err
	}
	stream := handleStreamResolution(resolution)

	download := downloader.New(downloader.Options{
		OutputName:   outputName,
		OutputPath:   outputPath,
		Stream:       stream,
		RetryTimes:   25,
		MultiThread:  true,
		ThreadNumber: 50,
	})
	log.Printf("the output is %v", outputPath)
	err = download.Download(data[0])
	if err != nil {
		log.Println("cleaning up, deleting folder...")
		if err := deleteContents(outputPath); err != nil {
			return "", err
		}
	   return "", err
	}
	fmt.Println("this is final output path", outputPath + outputName + ".*")
	files, err := filepath.Glob(outputPath + "/"+outputName + ".*")
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", errors.New("no video file found")
	}
	videoPath := files[0]
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return "", err
	}
	if fileInfo.Size() == 0 {
		return "", errors.New("downloaded video file is empty")
	}
	return videoPath, nil
}
