package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kingmariano/omnicron/config"
	ware "github.com/kingmariano/omnicron/middleware"
	"github.com/kingmariano/omnicron/packages/gpt"
	"github.com/kingmariano/omnicron/packages/grok"
	"github.com/kingmariano/omnicron/packages/musicdownloader"
	"github.com/kingmariano/omnicron/packages/replicate/generateimages"
	"github.com/kingmariano/omnicron/packages/replicate/generatemusic"
	"github.com/kingmariano/omnicron/packages/replicate/generatevideos"
	"github.com/kingmariano/omnicron/packages/replicate/imageupscale"
	"github.com/kingmariano/omnicron/packages/replicate/stt"
	"github.com/kingmariano/omnicron/packages/replicate/tts"
	"github.com/kingmariano/omnicron/packages/videodownloader"
	"github.com/kingmariano/omnicron/utils"
)

func setupRouter(t *testing.T) (*chi.Mux, *config.APIConfig) {
	apiKey, grokAPIKey, replicateAPIKey, cloudinaryURL, port, youtubeDeveloperKey, fastAPIPrivateURL, err := utils.LoadEnv("../.env")

	if err != nil {
		// if the given environment path is  not set. Get the variables from the root environment path
		apiKey = os.Getenv("API_KEY")
		grokAPIKey = os.Getenv("GROK_API_KEY")
		port = os.Getenv("PORT")
		replicateAPIKey = os.Getenv("REPLICATE_API_TOKEN")
		cloudinaryURL = os.Getenv("CLOUDINARY_URL")
		youtubeDeveloperKey = os.Getenv("YOUTUBE_DEVELOPER_KEY")
		fastAPIPrivateURL = os.Getenv("FAST_API_PRIVATE_URL")
	}

	if apiKey == "" || grokAPIKey == "" || replicateAPIKey == "" || cloudinaryURL == "" || port == "" || youtubeDeveloperKey == "" || fastAPIPrivateURL == "" {
		t.Fatal("unable to get API key or port from environment variables")
	}

	cfg := &config.APIConfig{
		APIKey:              apiKey,
		GrokAPIKey:          grokAPIKey,
		ReplicateAPIKey:     replicateAPIKey,
		CloudinaryURL:       cloudinaryURL,
		YoutubeDeveloperKey: youtubeDeveloperKey,
		Port:                port,
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", utils.HandleReadiness())
	v1Router.Post("/grok/chatcompletion", ware.MiddleWareAuth(grok.ChatCompletion, cfg))
	v1Router.Post("/grok/transcription", ware.MiddleWareAuth(grok.Transcription, cfg)) // deprecated
	v1Router.Post("/replicate/imagegeneration", ware.MiddleWareAuth(generateimages.ImageGeneration, cfg))
	v1Router.Post("/replicate/imageupscale", ware.MiddleWareAuth(imageupscale.ImageUpscale, cfg))
	v1Router.Post("/replicate/videogeneration", ware.MiddleWareAuth(generatevideos.VideoGeneration, cfg))
	v1Router.Post("/replicate/tts", ware.MiddleWareAuth(tts.TTS, cfg))
	v1Router.Post("/replicate/stt", ware.MiddleWareAuth(stt.STT, cfg))
	v1Router.Post("/replicate/musicgeneration", ware.MiddleWareAuth(generatemusic.MusicGen, cfg))
	v1Router.Post("/downloadvideo", ware.MiddleWareAuth(videodownloader.DownloadVideo, cfg))
	v1Router.Post("/downloadvideo", ware.MiddleWareAuth(musicdownloader.DownloadMusic, cfg))
	v1Router.Post("/chatgpt", ware.MiddleWareAuth(gpt.ChatCompletion, cfg))
	router.Mount("/api/v1", v1Router)

	return router, cfg
}

func TestReadiness(t *testing.T) {
	router, _ := setupRouter(t)

	req, err := http.NewRequest("GET", "/api/v1/readiness", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := http.StatusText(http.StatusOK)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
