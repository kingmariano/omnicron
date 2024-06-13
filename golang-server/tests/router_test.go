package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	ware "github.com/charlesozo/omnicron-backendsever/golang-server/middleware"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/grok"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/generateimages"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/generatemusic"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/generatevideos"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/imageupscale"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/stt"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/tts"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/videodownloader"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setupRouter(t *testing.T) (*chi.Mux, *config.ApiConfig) {
	// apiKey, grokApiKey, replicateApiKey, cloudinaryURL, port, err := utils.LoadEnv("../../.env")

	// if err != nil {
	// 	t.Fatal(err)
	// }

	apiKey := os.Getenv("API_KEY")
	grokApiKey := os.Getenv("GROK_API_KEY")
	port := os.Getenv("PORT")
	replicateApiKey := os.Getenv("REPLICATE_API_TOKEN")
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")

	if apiKey == "" || grokApiKey == "" || replicateApiKey == "" || cloudinaryURL == "" || port == "" {
		t.Fatal("unable to get API key or port from environment variables")
	}

	cfg := &config.ApiConfig{
		ApiKey:          apiKey,
		GrokApiKey:      grokApiKey,
		ReplicateAPIKey: replicateApiKey,
		CloudinaryUrl:   cloudinaryURL,
		Port:            port,
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
	v1Router.Post("/downloadvideo", ware.MiddleWareAuth(videodownloader.Download, cfg))
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

func TestChatCompletion(t *testing.T) {
	router, cfg := setupRouter(t)

	requestBody := `{"Model": "llama3-8b-8192", "Messages": [{"role": "user", "content": "Hello"}]}`
	req, err := http.NewRequest("POST", "/api/v1/grok/chatcompletion", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", cfg.ApiKey)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Further checks can be added based on the expected response
}

