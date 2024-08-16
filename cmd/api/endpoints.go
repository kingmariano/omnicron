package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/kingmariano/omnicron/config"
	ware "github.com/kingmariano/omnicron/middleware"
	"github.com/kingmariano/omnicron/packages/convert2mp3"
	"github.com/kingmariano/omnicron/packages/docgpt"
	"github.com/kingmariano/omnicron/packages/gpt"
	"github.com/kingmariano/omnicron/packages/grok"
	"github.com/kingmariano/omnicron/packages/image2text"
	"github.com/kingmariano/omnicron/packages/musicdownloader"
	"github.com/kingmariano/omnicron/packages/musicsearch"
	"github.com/kingmariano/omnicron/packages/replicate/generateimages"
	"github.com/kingmariano/omnicron/packages/replicate/generatemusic"
	"github.com/kingmariano/omnicron/packages/replicate/generatevideos"
	"github.com/kingmariano/omnicron/packages/replicate/imageupscale"
	"github.com/kingmariano/omnicron/packages/replicate/stt"
	"github.com/kingmariano/omnicron/packages/replicate/tts"
	"github.com/kingmariano/omnicron/packages/shazam"
	"github.com/kingmariano/omnicron/packages/videodownloader"
	"github.com/kingmariano/omnicron/packages/youtubesummarize"
	"github.com/kingmariano/omnicron/utils"
)

func callEndpoints(v1Router *chi.Mux, cfg *config.APIConfig){
	v1Router.Get("/readiness", utils.HandleReadiness())
	v1Router.Post("/groq/chatcompletion", ware.MiddleWareAuth(grok.ChatCompletion, cfg))
	v1Router.Post("/groq/transcription", ware.MiddleWareAuth(grok.Transcription, cfg)) // deprecated
	v1Router.Post("/replicate/imagegeneration", ware.MiddleWareAuth(generateimages.ImageGeneration, cfg))
	v1Router.Post("/replicate/imageupscale", ware.MiddleWareAuth(imageupscale.ImageUpscale, cfg))
	v1Router.Post("/replicate/videogeneration", ware.MiddleWareAuth(generatevideos.VideoGeneration, cfg))
	v1Router.Post("/replicate/tts", ware.MiddleWareAuth(tts.TTS, cfg))
	v1Router.Post("/replicate/stt", ware.MiddleWareAuth(stt.STT, cfg))
	v1Router.Post("/replicate/musicgeneration", ware.MiddleWareAuth(generatemusic.MusicGen, cfg))
	v1Router.Post("/downloadvideo", ware.MiddleWareAuth(videodownloader.DownloadVideo, cfg))
	v1Router.Post("/convert2mp3", ware.MiddleWareAuth(convert2mp3.ConvertToMp3, cfg))
	v1Router.Post("/downloadmusic", ware.MiddleWareAuth(musicdownloader.DownloadMusic, cfg))
	v1Router.Post("/gpt4free", ware.MiddleWareAuth(gpt.ChatCompletion, cfg))
	v1Router.Post("/shazam", ware.MiddleWareAuth(shazam.Shazam, cfg))
	v1Router.Post("/musicsearch", ware.MiddleWareAuth(musicsearch.MusicSearch, cfg))
	v1Router.Post("/youtubesummarization", ware.MiddleWareAuth(youtubesummarize.YoutubeSummarization, cfg))
	v1Router.Post("/image2text", ware.MiddleWareAuth(image2text.Image2text, cfg))
	v1Router.Post("/docgpt", ware.MiddleWareAuth(docgpt.DocGPT, cfg))
}