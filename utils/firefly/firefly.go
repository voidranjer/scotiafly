package firefly

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/voidranjer/scotiafly/utils"
)

const URL = "http://localhost/api/v1/transactions"

type TransactionPayload struct {
	TransactionType string
	Description     string
	CategoryName    string
	Amount          float32
	Date            string
	ExternalID      string
	AccountName     string
}

func PostTransaction(payload TransactionPayload) (success bool) {
	var sourceOrDest string

	switch payload.TransactionType {
	case "withdrawal":
		sourceOrDest = "source_name"
	case "deposit":
		sourceOrDest = "destination_name"
	default:
		log.Fatal("PostTransaction: Invalid transaction type!")
	}

	body := map[string]any{
		"error_if_duplicate_hash": true,
		"transactions": []map[string]any{
			{
				"type":          payload.TransactionType,
				"description":   payload.Description,
				"category_name": payload.CategoryName,
				"amount":        payload.Amount,
				"date":          payload.Date,
				"external_id":   payload.ExternalID,
				sourceOrDest:    payload.AccountName,
			},
		},
	}

	marshaledBody, err := json.Marshal(body)
	utils.HandleError(err)

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(marshaledBody))
	utils.HandleError(err)

	utils.AttachHeadersToRequest("headers_firefly.txt", req)

	client := &http.Client{}

	response, err := client.Do(req)
	utils.HandleError(err)

	if response.StatusCode != http.StatusOK {
		log.Println("Error:", response.Status)

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
		} else {
			log.Println("Error response body:", string(bodyBytes))
		}

		return false
	}

	// log.Println("Posted: ", )
	return true
}
