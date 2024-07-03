package musicsearch

import (
	"encoding/json"
	"fmt"
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

type MusicSearchRequest struct {
	Song  string `json:"song"`
	Limit int    `json:"limit,omitempty"`
	Proxy string `json:"proxy,omitempty"`
}
type MusicSearchResponse struct {
	Tracks struct {
		Hits []struct {
			Type    string `json:"type"`
			Key     string `json:"key"`
			Heading struct {
				Title    string `json:"title"`
				Subtitle string `json:"subtitle"`
			} `json:"heading"`
			Images struct {
				Default string `json:"default"`
				Blurred string `json:"blurred"`
				Play    string `json:"play"`
			} `json:"images"`
			Stores struct {
				Apple struct {
					Actions []struct {
						Type string `json:"type"`
						URI  string `json:"uri"`
					} `json:"actions"`
					Explicit    bool   `json:"explicit"`
					Previewurl  string `json:"previewurl"`
					Coverarturl string `json:"coverarturl"`
					Trackid     string `json:"trackid"`
					Productid   string `json:"productid"`
				} `json:"apple"`
			} `json:"stores"`
			Streams struct {
			} `json:"streams"`
			Artists []struct {
				Alias  string `json:"alias"`
				ID     string `json:"id"`
				Adamid string `json:"adamid"`
			} `json:"artists"`
			Share struct {
				Subject  string `json:"subject"`
				Text     string `json:"text"`
				Href     string `json:"href"`
				Image    string `json:"image"`
				Twitter  string `json:"twitter"`
				HTML     string `json:"html"`
				Avatar   string `json:"avatar"`
				Snapchat string `json:"snapchat"`
			} `json:"share"`
			Alias   string `json:"alias"`
			URL     string `json:"url"`
			Actions []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"actions"`
			Urlparams struct {
				Tracktitle  string `json:"{tracktitle}"`
				Trackartist string `json:"{trackartist}"`
			} `json:"urlparams"`
		} `json:"hits"`
		Next string `json:"next"`
	} `json:"tracks"`
}

type ErrorResponse struct {
	Detail string `json:"detail"`
}

func MusicSearch(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	decode := json.NewDecoder(r.Body)
	musicSearchParams := MusicSearchRequest{}
	err := decode.Decode(&musicSearchParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling json, %v", err))
		return
	}
	response, err := CallMusicSearchFastAPI(musicSearchParams, cfg.ApiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}
