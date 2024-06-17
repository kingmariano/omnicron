package video2audio
import (
	"os"
	"io"
	"net/http"
)

// downloadFile downloads a file from the given URL and saves it to the specified path.
func downloadFile(url string, dest string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    out, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}