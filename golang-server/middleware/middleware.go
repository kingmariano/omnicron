package middleware
import (
	"fmt"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
	"github.com/charlesozo/omnicron-backendsever/golang-server/internal/auth"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"net/http"
)

type authHandler func(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig)

func MiddleWareAuth(handler authHandler, cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetHeaderToken(r.Header)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprint(err))
			return
		}
		if token != cfg.ApiKey {
			utils.RespondWithError(w, http.StatusUnauthorized, "Api Key is invalid")
			return
		}

		handler(w, r, cfg)
	}
}
