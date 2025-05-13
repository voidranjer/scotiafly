package utils

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
)

func ParseHeadersFromFile(filePath string) map[string]string {
	result := make(map[string]string)

	file, err := os.Open(filePath)
	HandleError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue // or log warning
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		result[key] = value
	}

	HandleError(scanner.Err())

	return result
}

func AttachHeadersToRequest(filePath string, request *http.Request) {
	headers := ParseHeadersFromFile(filePath)

	for key, value := range headers {
		request.Header.Add(key, value)
	}
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
