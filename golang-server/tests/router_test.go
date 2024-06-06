package tests

import (
	"bytes"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	ware "github.com/charlesozo/omnicron-backendsever/golang-server/middleware"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/grok"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/generateimages"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/generatevideos"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/imageupscale"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
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
	cloudinaryURL  := os.Getenv("CLOUDINARY_URL")

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
	v1Router.Post("/replicate/imagegeneration", ware.MiddleWareAuth(generateimages.ImageGeneration, cfg))
	v1Router.Post("/replicate/imageupscale", ware.MiddleWareAuth(imageupscale.ImageUpscale, cfg))
	v1Router.Post("/replicate/videogeneration", ware.MiddleWareAuth(generatevideos.VideoGeneration, cfg))
	v1Router.Post("/grok/transcription", ware.MiddleWareAuth(grok.Transcription, cfg)) // deprecated
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
func TestLowImageGeneration(t *testing.T) {
	router, cfg := setupRouter(t)
	requestBody := `{"prompt": "self-portrait of a woman, lightning in the background"}`
	baseURL := "/api/v1/replicate/imagegeneration"
	params := url.Values{}
	params.Add("model", "bytedance/sdxl-lightning-4step")
	url := baseURL + "?" + params.Encode()
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
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

func TestHighImageGeneration(t *testing.T) {
	router, cfg := setupRouter(t)

	// Create a buffer to hold the form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add form fields
	w.WriteField("prompt", "A detailed portrait of a cyberpunk woman with neon blue hair, intricate facial tattoos, and a metallic outfit.")
	// _ = w.WriteField("width", "512")
	// _ = w.WriteField("height", "512")
	// _ = w.WriteField("num_outputs", "1")
	// _ = w.WriteField("guidance_scale", "7.5")
	// _ = w.WriteField("num_inference_steps", "50")
	// _ = w.WriteField("apply_watermark", "true")
	// _ = w.WriteField("scheduler", "ddim")

	// Optional: Add image file
	imagePath := "../../assets/images/test_image1.png"
	file, err := os.Open(imagePath)
	if err == nil {
		fw, err := w.CreateFormFile("image", "image.png")
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(fw, file)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}

	// Close the multipart writer to flush the buffer
	w.Close()

	baseURL := "/api/v1/replicate/imagegeneration"
	params := url.Values{}
	params.Add("model", "lorenzomarines/astra")
	url := baseURL + "?" + params.Encode()
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Api-Key", cfg.ApiKey)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Log(rr.Body)
	}

	// Further checks can be added based on the expected response
}

func TestImageUpscale(t *testing.T) {
	router, cfg := setupRouter(t)

	// Create a buffer to hold the form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	imagePath := "../../assets/images/test_image1.png"
	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	fw, err := w.CreateFormFile("image", "image.png")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	// Close the multipart writer to flush the buffer
	w.Close()
	baseURL := "/api/v1/replicate/imageupscale"
	params := url.Values{}
	params.Add("model", "nightmareai/real-esrgan")
	url := baseURL + "?" + params.Encode()
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Api-Key", cfg.ApiKey)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Log(rr.Body)
	}

	// Further checks can be added based on the expected response
}

func TestVideoGeneration(t *testing.T) {
	router, cfg := setupRouter(t)

	// Create a buffer to hold the form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	w.WriteField("prompt", "A sunrise over a calm ocean, waves gently lapping the shore, seagulls flying overhead.")
	// Close the multipart writer to flush the buffer
	w.Close()
	baseURL := "/api/v1/replicate/videogeneration"
	params := url.Values{}
	params.Add("model", "anotherjesse/zeroscope-v2-xl")
	url := baseURL + "?" + params.Encode()
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Api-Key", cfg.ApiKey)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Log(rr.Body)
	}

	// Further checks can be added based on the expected response
}
