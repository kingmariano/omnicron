package tests

import (
	"context"
	"testing"

	"github.com/charlesozo/omnicron-backendsever/golang-server/storage"
)

func TestUpload(t *testing.T) {
	_, cfg := setupRouter(t)
	filePath := "../../assets/images/test_image1.png"
	_, err := storage.HandleFileUpload(context.Background(), filePath, cfg)
	if err != nil {
		t.Fatal(err)
	}
}
