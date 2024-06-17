package utils
import (
	"os"
	"log"
)

func CreateFolder(outputPath string){
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		err := os.Mkdir(outputPath, os.ModePerm)
		if err != nil {
			log.Printf("error creating directory %v", err)
			return
			// utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating directory, %v", err))
		}
		defer os.RemoveAll(outputPath) // Ensure the folder is deleted after task completion
	}
	log.Print("file already exists")
	return
}