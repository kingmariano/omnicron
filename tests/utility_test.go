package tests

import (
	"github.com/kingmariano/omnicron/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateAndDeleteTemporaryFolderMultipleCalls(t *testing.T) {
	// Test case: Create multiple temporary directories and delete them after the function has completed
	for i := 0; i <= 5; i++ {
		folderPath, err := utils.CreateUniqueFolder(utils.BasePath)
		if err != nil {
			t.Fatalf("unexpected error creating  directory: %v", err)
		}
		err = utils.DeleteFolder(folderPath)
		if err != nil {
			t.Fatalf("unexpected error deleting temporary directory: %v", err)
		}

		// Check if the temporary directory has been deleted
		_, err = os.Stat(folderPath)
		if !os.IsNotExist(err) {
			t.Errorf("temporary directory still exists after deletion: %s", folderPath)
		}
	}
}

func TestDownloadFileURL(t *testing.T) {
	testCases := []struct {
		url         string
		shouldExist bool
	}{
		{"https://res.cloudinary.com/djagytapi/video/upload/v1718161686/zsq3vzbjbtg7eqrpi2ui.mp4", true},
		{"https://res.cloudinary.com/djagytapi/video/upload/v1718249303/nm1eiabpo8zwo6hz6uxh.mp4", true},
		{"https://www.example.com/nonexistent.mp3", false},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			folderPath, err := utils.CreateUniqueFolder(utils.BasePath)
			if err != nil {
				t.Fatalf("unexpected error creating  directory: %v", err)
			}
			downloadedFile, err := utils.DownloadFileURL(tc.url, folderPath)
			if tc.shouldExist {
				assert.NoError(t, err)
				assert.FileExists(t, downloadedFile)

				// Clean up
				err := utils.DeleteFolder(folderPath)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, "", downloadedFile)
			}
			//final cleanup of all directories
			err = utils.DeleteFolder(folderPath)
			assert.NoError(t, err)
		})
	}
}

func TestConvertToMP3(t *testing.T) {
	filePath := "../assets/audios/sample1.aiff"
	file, err := utils.ConvertFileToMP3(filePath)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, file)
	assert.NoError(t, os.Remove(file))
}

func TestConvertReaderToMP3(t *testing.T) {
	tests := []struct {
		FilePath        string
		ExpectedOutcome error
	}{
		{
			FilePath:        "../assets/audios/sample1.aiff",
			ExpectedOutcome: nil,
		},
		{
			FilePath:        "../assets/videos/sample.mp4",
			ExpectedOutcome: nil,
		},
		{
			FilePath:        "../assets/documents/sample1.pdf",
			ExpectedOutcome: utils.UnsupportedFileFormat,
		},
	}
	for _, test := range tests {
		t.Run(test.FilePath, func(t *testing.T) {
			folderPath, err := utils.CreateUniqueFolder(utils.BasePath)
			if err != nil {
				t.Fatalf("unexpected error creating  directory: %v", err)
			}
			videoFile, err := os.Open(test.FilePath)
			if err != nil {
				//remove the existing folder created
				err = utils.DeleteFolder(folderPath)
				assert.NoError(t, err)
				t.Fatalf("unexpected error opening file: %v", err)
			}
			// Convert the reader to MP3 and check the outcome
			outputFileName, err := utils.ConvertReaderToMP3(videoFile, folderPath)
			assert.Equal(t, test.ExpectedOutcome, err)
			if err == nil {
				assert.NotEmpty(t, outputFileName)
			}
			//clean up the directory
			err = utils.DeleteFolder(folderPath)
			assert.NoError(t, err)
		})
	}

}

// TestConvertURLToMP3 tests the downloading and converting of a URL to MP3
func TestConvertURLToMP3(t *testing.T) {
	testURLs := []string{
		"https://res.cloudinary.com/djagytapi/video/upload/v1718161686/zsq3vzbjbtg7eqrpi2ui.mp4",
		"https://res.cloudinary.com/djagytapi/video/upload/v1718249303/nm1eiabpo8zwo6hz6uxh.mp4",
	}

	for _, url := range testURLs {
		t.Run(url, func(t *testing.T) {
			folderPath, err := utils.CreateUniqueFolder(utils.BasePath)
			assert.NoError(t, err)
			downloadedFile, err := utils.DownloadFileURL(url, folderPath)
			assert.NoError(t, err)

			file, err := utils.ConvertFileToMP3(downloadedFile)
			assert.NoError(t, err)
			assert.NotEmpty(t, file)

			// Clean up
			err = utils.DeleteFolder(folderPath)
			assert.NoError(t, err)
		})
	}
}
