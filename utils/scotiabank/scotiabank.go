package scotiabank

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/voidranjer/scotiafly/utils"
)

func MakeHeaderedRequest(url string, folderName string) (result Response) {
	req, err := http.NewRequest("GET", url, nil)
	utils.HandleError(err)
	utils.AttachHeadersToRequest("headers_scotia.txt", req)

	client := &http.Client{}

	response, err := client.Do(req)
	utils.HandleError(err)

	reader := response.Body
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		utils.HandleError(err)
	}

	defer reader.Close()

	body, err := io.ReadAll(reader)
	utils.HandleError(err)

	// write to debug file
	var f map[string]any
	json.Unmarshal(body, &f)
	b, err := json.MarshalIndent(f, "", "  ")
	utils.HandleError(err)
	timestamp := time.Now().Format("2006-01-02_15-04-05.000")
	os.MkdirAll(folderName, 0755)
	filename := folderName + "/output_" + timestamp + ".json"
	os.WriteFile(filename, b, 0644)

	err = json.Unmarshal(body, &result)
	utils.HandleError(err)

	return result
}
