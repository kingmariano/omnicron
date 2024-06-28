package musicdownloader

import (
	"context"
	"testing"
)

func TestSearchMusicOnYoutube(t *testing.T) {
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
		},
		{
			name: "Test11",
			args: args{
				songTitle:          "Sorry - Nothing But Thieves",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=p9_BsjMi4bM",
			},
		},
		{
			name: "Test12",
			args: args{
				songTitle:          "Sorry - Justin Bieber",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=fRh_vgS2dFE",
			},
		},
		{
			name: "Test13",
			args: args{
				songTitle:          "El Anciano y el Niño (Acústica) - Cheo Gallego",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=4AMdEQpLF6Q",
			},
		},
		{
			name: "Test14",
			args: args{
				songTitle:          "One Dance (feat. Wizkid & Kyla) - Drake",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=ki0Ocze98U8",
			},
		},
		{
			name: "Test15",
			args: args{
				songTitle:          "Brothers - Lil Tjay",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=ifZzXeNt3L4",
			},
		},
		{
			name: "Test16",
			args: args{
				songTitle:          "One Dance / Hasta El Amanecer (Mashup) - Alex Aiono",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=Kuz3DUNZaC8",
			},
		},
		{
			name: "Test17",
			args: args{
				songTitle:          "Shake It Up - Trippie Redd",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=vDw-y5w9sG0",
			},
		},
		{
			name: "Test18",
			args: args{
				songTitle:          "Shake It Up - Selena Gomez",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=VVchUL0ZK_4",
			},
		},
		{
			name: "Test19",
			args: args{
				songTitle:          "Mood Swings (feat. Lil Tjay) - Pop Smoke",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=mM8ostx0Ub8",
			},
		},
		{
			name: "Test20",
			args: args{
				songTitle:          "Mood Swings (feat. Lil Tjay & Summer Walker) [Remix] - Pop Smoke",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=UqO4dd3Ea9Q",
			},
		},
		{
			name: "Test21",
			args: args{
				songTitle:          "shut up - Ariana Grande",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=9MogWz-LHXI",
			},
		},
		{
			name: "Test22",
			args: args{
				songTitle:          "Shut Up - Stormzy",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=RqQGUJK7Na4",
			},
		},
		{
			name: "Test23",
			args: args{
				songTitle:          "Shut Up - Madness",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=Ul_Kotlqi3k",
			},
		},
		{
			name: "Test24",
			args: args{
				songTitle:          "Hightension - Shallipopi",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=HPtQ_YMYmus",
			},
		},
		{
			name: "Test25",
			args: args{
				songTitle:          "Crazy - Seal",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=4Fc67yQsPqQ",
			},
		},
		{
			name: "Test26",
			args: args{
				songTitle:          "Crazy - Gnarls Barkley",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=-N4jf6rtyuw",
			},
		},
		{
			name: "Test27",
			args: args{
				songTitle:          "Crazy - Aerosmith",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=NMNgbISmF4I",
			},
		},
		{
			name: "Test28",
			args: args{
				songTitle:          "Crazy - Daniela Andrade",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=fzxag7U3Snk",
			},
		},
		{
			name: "Test29",
			args: args{
				songTitle:          "No Better Love (feat. Rell) - Young Gunz, featuring Rell",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=3kpzbqV67ZU",
			},
		},
		{
			name: "Test30",
			args: args{
				songTitle:          "Count Me In - Cast - Liv and Maddie",
				expectedYoutubeURL: "https://www.youtube.com/watch?v=EHl4Ht1Xbe0",
			},
		},
	}
	ctx := context.TODO()
	maxResult := int64(1)
	youtubeApTestingKey := "AIzaSyCSV2NjjOxnn0GRkl7FaFPJ5evYTr1Yc_E"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlList, err := searchMusicOnYouTube(ctx, tt.args.songTitle, maxResult, youtubeApTestingKey, "")
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
