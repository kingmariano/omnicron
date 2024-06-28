package musicdownloader

// import (
// 	"context"
// 	"testing"
// )

// func TestSearchMusicOnYoutube(t *testing.T) {
// 	type args struct {
// 		songTitle          string
// 		expectedYoutubeURL string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			name: "Test1",
// 			args: args{
// 				songTitle:          "Brothers (feat. Anthony Thomas) - Krept & Konan",
// 				expectedYoutubeURL: "https://www.youtube.com/watch?v=siGrcXUiuko",
// 			},
// 		},
// 		{
// 			name: "Test2",
// 			args: args{
// 				songTitle:          "Brothers - Lil Tjay",
// 				expectedYoutubeURL: "https://www.youtube.com/watch?v=ifZzXeNt3L4",
// 			},
// 		},
// 		{
// 			name: "Test3",
// 			args: args{
// 				songTitle:          "La Isla Bonita - Madonna",
// 				expectedYoutubeURL: "https://www.youtube.com/watch?v=zpzdgmqIHOQ",
// 			},
// 		},
// 		{
// 			name: "Test4",
// 			args: args{
// 				songTitle:          "Faded - Conor Maynard",
// 				expectedYoutubeURL: "https://www.youtube.com/watch?v=IgCphQCkHSk",
// 			},
// 		},
// 		{
// 			name: "Test5",
// 			args: args{
// 				songTitle:          "Faded (Restrung) - Alan Walker",
// 				expectedYoutubeURL: "https://www.youtube.com/watch?v=bDmzGLrdjxQ",
// 			},
// 		},
// 		{
// 			name: "Test6",
// 			args: args{
// 				songTitle:          "Watermelon - Justine Clarke",
// 				expectedYoutubeURL: "https://www.youtube.com/watch?v=yhujKmHDEEc",
// 			},
// 		},
// 		{
// 			name: "Test7",
// 			args: args{
// 				songTitle:          "Chlorine - Title Fight",
// 				expectedYoutubeURL: "https://www.youtube.com/watch?v=YJRGAcu2fb8",
// 			},
// 		},
// 		{
// 			name: "Test8",
// 			args: args{
// 				songTitle:          "Watermelon Sugar - Harry Styles",
// 				expectedYoutubeURL: "https://www.youtube.com/watch?v=E07s5ZYygMg",
// 			},
// 		},
// 		{
// 			name: "Test9",
// 			args: args{
// 				songTitle:          "Chlorine - twenty one pilots",
// 				expectedYoutubeURL: "https://www.youtube.com/watch?v=eJnQBXmZ7Ek",
// 			},
// 		},
// 		{
// 			name: "Test10",
// 			args: args{
// 				songTitle:          "I Love You, I Hate You - Little Simz",
// 				expectedYoutubeURL: "https://www.youtube.com/watch?v=GSOFGFagOsY",
// 			},
// 		}}
// 	ctx := context.TODO()
// 	maxResult := int64(1)
// 	youtubeApTestingKey := "AIzaSyCSV2NjjOxnn0GRkl7FaFPJ5evYTr1Yc_E"

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			urlList, err := searchMusicOnYouTube(ctx, tt.args.songTitle, maxResult, youtubeApTestingKey, "")
// 			if err != nil {
// 				t.Errorf("searchMusicOnYouTube() error = %v", err)
// 				return
// 			}
// 			if urlList[0] != tt.args.expectedYoutubeURL {
// 				t.Errorf("searchMusicOnYouTube() = %v, want %v", urlList[0], tt.args.expectedYoutubeURL)
// 			}
// 		})
// 	}
// }
