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
package tests

import (
	"context"
	"testing"

	"github.com/kingmariano/omnicron/packages/musicdownloader"
)

func TestSearchMusicOnYoutube(t *testing.T) {
	ctx := context.Background()
	maxResult := int64(1)
	_, cfg := setupRouter(t)
	type args struct {
		songTitle          string
		expectedYoutubeURL string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test1",
			args: args{
				songTitle:          "Brothers (feat. Anthony Thomas) - Krept & Konan",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=siGrcXUiuko",
			},
		},
		{
			name: "Test2",
			args: args{
				songTitle:          "Brothers - Lil Tjay",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=ifZzXeNt3L4",
			},
		},
		{
			name: "Test3",
			args: args{
				songTitle:          "La Isla Bonita - Madonna",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=zpzdgmqIHOQ",
			},
		},
		{
			name: "Test4",
			args: args{
				songTitle:          "Faded - Conor Maynard",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=IgCphQCkHSk",
			},
		},
		{
			name: "Test5",
			args: args{
				songTitle:          "Faded (Restrung) - Alan Walker",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=bDmzGLrdjxQ",
			},
		},
		{
			name: "Test6",
			args: args{
				songTitle:          "Watermelon - Justine Clarke",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=yhujKmHDEEc",
			},
		},
		{
			name: "Test7",
			args: args{
				songTitle:          "Chlorine - Title Fight",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=YJRGAcu2fb8",
			},
		},
		{
			name: "Test8",
			args: args{
				songTitle:          "Watermelon Sugar - Harry Styles",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=E07s5ZYygMg",
			},
		},
		{
			name: "Test9",
			args: args{
				songTitle:          "Chlorine - twenty one pilots",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=eJnQBXmZ7Ek",
			},
		},
		{
			name: "Test10",
			args: args{
				songTitle:          "I Love You, I Hate You - Little Simz",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=GSOFGFagOsY",
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlList, err := musicdownloader.SearchMusicOnYouTube(ctx, tt.args.songTitle, maxResult, cfg.YoutubeDeveloperKey, "")
			if err != nil {
				t.Errorf("searchMusicOnYouTube() error = %v", err)
				return
			}
			if urlList[0] != tt.args.expectedYoutubeURL {
				t.Errorf("searchMusicOnYouTube() = %v, want %v", urlList[0], tt.args.expectedYoutubeURL)
			}
		})
	}
}
