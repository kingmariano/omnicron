package api
import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kingmariano/omnicron/config"
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

func StartServer() {
	// Load environment variables
	apiKey, grokAPIKey, replicateAPIKey, cloudinaryURL, port, fastAPIBaseURL, err := utils.LoadEnv()
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
	callEndpoints(v1Router, cfg)
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
