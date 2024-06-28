package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatCompletion(t *testing.T) {
	router, cfg := setupRouter(t)

	requestBody := `{"Model": "llama3-70b-8192", "Messages": [{"role": "user", "content": "Hello"}]}`
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
		t.Log(rr.Body)
	}

	// Further checks can be added based on the expected response
}
