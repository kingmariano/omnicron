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
package videodownloader

import (
	"github.com/kingmariano/omnicron/utils"
	"testing"
)

func TestDownloadVideoData(t *testing.T) {
	type args struct {
		url        string
		outputName string
		resolution string
	}

	tests := []struct {
		name string
		args args
	}{

		// {
		// 	name: "Download 1080p video",
		// 	args: args{
		// 		url:        "https://www.youtube.com/shorts/WO7wT-FX2mA",
		// 		outputName: "test_video_1080p",
		// 		resolution: "",
		// 	},
		// },
		// {
		// 	name: "Download 720p video",
		// 	args: args{
		// 		url:        "https://youtu.be/ZT0yQgUIZho",
		// 		outputName: "test_video_720p",
		// 		resolution: "720p",
		// 	},
		// },

		{
			name: "Download 240p video",
			args: args{
				url:        "https://youtu.be/ZT0yQgUIZho",
				outputName: "test_video_240p",
				resolution: "360p",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the output directory before running tests
			folderPath, err := utils.CreateUniqueFolder(utils.BasePath)
			if err != nil {
				t.Errorf("Failed to create unique folder: %v", err)
				return
			}
			_, err = DownloadVideoData(tt.args.url, tt.args.outputName, folderPath, tt.args.resolution)
			if err != nil {
				t.Errorf("DownloadVideoData() error = %v", err)
				return
			}
			//delete folder after downloading
			err = utils.DeleteFolder(folderPath)
			if err != nil {
				t.Errorf("Failed to cleanup: %v", err)
			}
			t.Logf("removing folder: %v", folderPath)
		})

	}

}
