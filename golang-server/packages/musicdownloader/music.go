package musicdownloader
import (
	"fmt"
	"net/http"
	"sync"
	"log"
	"context"
	"encoding/json"
	"github.com/charlesozo/omnicron-backendsever/golang-server/utils"
	"github.com/charlesozo/omnicron-backendsever/golang-server/config"
)
type SongRequest struct {
	Songs      []string `json:"songs"`

}
// SongResponse represents the structure of the response
type SongResponse map[string][]string
var maxLength int64 //specifies the maxmium length of data returned from the youtube sdk
func DownloadMusic(w http.ResponseWriter, r *http.Request, cfg *config.ApiConfig){
	log.Print("started downloading music process")
	ctx := context.Background()
	decode := json.NewDecoder(r.Body)
	params := SongRequest{}
	err := decode.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling json, %v", err))
		return
	}
	var wg sync.WaitGroup
	
	resultChan := make(chan map[string][]string, len(params.Songs))
	maxLength = 1
	for _, song := range params.Songs{
		sng := song
		wg.Add(1)
		go searchYouTube(ctx, sng, maxLength, cfg.YoutubeDeveloperKey, &wg, resultChan)
	}
	// Close the result channel once all searches are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()
    songResponse := make(SongResponse)
	for result := range resultChan{
		for k, v := range result {
			songResponse[k] = v
		}
	}
	utils.RespondWithJSON(w, http.StatusOK, songResponse)
}