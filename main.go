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

package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

// starts the python fastAPI server
func startFastAPIServer() *exec.Cmd {
	cmd := exec.Command("python", "./python/main.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start FastAPI server: %v", err)
	}

	// Give the server a few seconds to start
	time.Sleep(5 * time.Second)

	return cmd
}

func main() {
	// Load environment variables
	apiKey, grokAPIKey, replicateAPIKey, cloudinaryURL, port, fastAPIBaseURL, err := utils.LoadEnv(".env")
	if err != nil {
		log.Fatalf("error opening .env %v ", err)
	}

	// Start the FastAPI server
	cmd := startFastAPIServer()
	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("Failed to kill process: %v", err)
		}
	}()

	cfg := &config.APIConfig{
		APIKey:          apiKey,
		GrokAPIKey:      grokAPIKey,
		ReplicateAPIKey: replicateAPIKey,
		CloudinaryURL:   cloudinaryURL,
		FASTAPIBaseURL:  fastAPIBaseURL,
		Port:            port,
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	v1Router := chi.NewRouter()
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
	router.Mount("/api/v1", v1Router)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Create a channel to listen for OS signals for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server started at port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	<-stop

	log.Println("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
