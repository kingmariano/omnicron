package main

import (
	"log"
	"net/http"
	"time"

	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	ware "github.com/charlesozo/omnicron-backendsever/golang-server/middleware"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/grok"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/generateimages"
	"github.com/charlesozo/omnicron-backendsever/golang-server/packages/replicate/imageupscale"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	apiKey, grokApiKey, replicateApiKey, cloudinaryURL, port, err := utils.LoadEnv("../.env")
	if err != nil {
		log.Fatal(err)
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
	v1Router.Post("/grok/transcription", ware.MiddleWareAuth(grok.Transcription, cfg)) // deprecated
	router.Mount("/api/v1", v1Router)
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server started at port %s", port)
	log.Fatal(server.ListenAndServe())
}
