package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, errmessage string) {
	if code < 499 {
		log.Printf("Responding with 5XX error: %s", errmessage)
	}
	type Errormsg struct {
		Error string `json:"error"`
	}
	RespondWithJSON(w, code, Errormsg{
		Error: errmessage,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
        log.Print(err)
		return
	}
}
