package shazam

import (
	"fmt"
	"github.com/kingmariano/omnicron/config"
	"github.com/kingmariano/omnicron/utils"
	"net/http"
)

type ShazamResponse struct {
	Matches []struct {
		ID            string  `json:"id"`
		Offset        float64 `json:"offset"`
		Timeskew      float64 `json:"timeskew"`
		Frequencyskew float64 `json:"frequencyskew"`
	} `json:"matches"`
	Location struct {
		Accuracy float64 `json:"accuracy"`
	} `json:"location"`
	Timestamp int    `json:"timestamp"`
	Timezone  string `json:"timezone"`
	Track     struct {
		Layout   string `json:"layout"`
		Type     string `json:"type"`
		Key      string `json:"key"`
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Images   struct {
			Background string `json:"background"`
			Coverart   string `json:"coverart"`
			Coverarthq string `json:"coverarthq"`
			Joecolor   string `json:"joecolor"`
		} `json:"images"`
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
		Hub struct {
			Type    string `json:"type"`
			Image   string `json:"image"`
			Actions []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				ID   string `json:"id,omitempty"`
				URI  string `json:"uri,omitempty"`
			} `json:"actions"`
			Options []struct {
				Caption string `json:"caption"`
				Actions []struct {
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"actions"`
				Beacondata struct {
					Type         string `json:"type"`
					Providername string `json:"providername"`
				} `json:"beacondata"`
				Image               string `json:"image"`
				Type                string `json:"type"`
				Listcaption         string `json:"listcaption"`
				Overflowimage       string `json:"overflowimage"`
				Colouroverflowimage bool   `json:"colouroverflowimage"`
				Providername        string `json:"providername"`
			} `json:"options"`
			Explicit    bool   `json:"explicit"`
			Displayname string `json:"displayname"`
		} `json:"hub"`
		Sections []struct {
			Type      string `json:"type"`
			Metapages []struct {
				Image   string `json:"image"`
				Caption string `json:"caption"`
			} `json:"metapages,omitempty"`
			Tabname  string `json:"tabname"`
			Metadata []struct {
				Title string `json:"title"`
				Text  string `json:"text"`
			} `json:"metadata,omitempty"`
			URL string `json:"url,omitempty"`
		} `json:"sections"`
		URL     string `json:"url"`
		Artists []struct {
			Alias  string `json:"alias"`
			ID     string `json:"id"`
			Adamid string `json:"adamid"`
		} `json:"artists"`
		Alias  string `json:"alias"`
		Isrc   string `json:"isrc"`
		Genres struct {
			Primary string `json:"primary"`
		} `json:"genres"`
		Urlparams struct {
			Tracktitle  string `json:"{tracktitle}"`
			Trackartist string `json:"{trackartist}"`
		} `json:"urlparams"`
		Highlightsurls struct {
		} `json:"highlightsurls"`
		Relatedtracksurl string `json:"relatedtracksurl"`
		Albumadamid      string `json:"albumadamid"`
		Trackadamid      string `json:"trackadamid"`
		Releasedate      string `json:"releasedate"`
	} `json:"track"`
	Tagid string `json:"tagid"`
}

type ErrorResponse struct {
	Detail string `json:"detail"`
}

func Shazam(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig) {
	// Parse the multipart form in the request
	err := r.ParseMultipartForm(30 << 20) // 30MB max memory
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing multipart form, %v", err))
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error retrieving the file, %v", err))
		return
	}
	defer file.Close()
	response, err := CallShazamFastAPI(file, fileHeader, cfg.ApiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error calling Shazam API, %v", err))
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}
