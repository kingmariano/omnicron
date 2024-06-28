package auth

import (
	"errors"
	"net/http"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func GetHeaderToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Api-Key")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	return authHeader, nil
}
