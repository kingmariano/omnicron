package videodownloader

import (
	_ "github.com/iawia002/lux/app"
	"github.com/iawia002/lux/downloader"
	"github.com/iawia002/lux/extractors"
	"log"
	"os"
"path/filepath"
)

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

func DownloadVideoData(url string, outputName string, outputPath string, resolution string) error {
	data, err := extractUrl(url)
	if err != nil {
		return err
	}
	
   stream := handleStreamResolution(resolution)

	download := downloader.New(downloader.Options{
		OutputName:   outputName,
		OutputPath:   outputPath,
		Stream:       stream,
		RetryTimes: 25,
		MultiThread:  true,
		ThreadNumber: 50,
	})
	log.Printf("the output is %v", outputPath)
	err = download.Download(data[0])
	if err != nil {
		log.Print("cleaning up, removing unnecessary files")
		if err := deleteContents(outputPath); err != nil {
			return err
		}
		return err
	}
	return nil
}
