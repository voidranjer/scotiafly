package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/voidranjer/scotiafly/utils"
	"github.com/voidranjer/scotiafly/utils/firefly"
	"github.com/voidranjer/scotiafly/utils/scotiabank"
)

const MAX_PAGINATION = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./scotiafly <path_to_config: [config_chequing.json | config_sceneplus.json]>")
		os.Exit(1)
	}

	dat, err := os.ReadFile(os.Args[1])
	utils.HandleError(err)
	var config Config
	json.Unmarshal(dat, &config)

	today := time.Now()
	todayFormatted := today.Format("2006-01-02")

	twoYearsAgo := today.AddDate(-2, 0, 0)
	twoYearsAgoFormatted := twoYearsAgo.Format("2006-01-02")

	fromDate := twoYearsAgoFormatted
	toDate := todayFormatted

	params := fmt.Sprintf("?fromDate=%s&toDate=%s&limit=%d", fromDate, toDate, MAX_PAGINATION)
	url := config.AccountUrl + params

	cursoredUrl := url

	for {
		response := scotiabank.MakeHeaderedRequest(cursoredUrl, config.FolderName)

		for _, transaction := range response.Transactions {
			var transactionType string
			switch transaction.TransactionType {
			case "CREDIT":
				transactionType = "deposit"
			case "DEBIT":
				transactionType = "withdrawal"
			}

			success := firefly.PostTransaction(firefly.TransactionPayload{
				TransactionType: transactionType,
				Description:     transaction.CleanDescription,
				CategoryName:    transaction.Category.Description,
				Amount:          transaction.TransactionAmount.Amount,
				Date:            transaction.TransactionDate,
				ExternalID:      transaction.Id,
				AccountName:     config.AccountName,
			})

			if !success {
				fmt.Println("Warning: Failed to post transaction ", transaction.CleanDescription)
			} else {
				fmt.Println("Successfully posted ", transaction.CleanDescription)
			}
		}

		if response.NextCursorKey == "" {
			break
		}

		cursoredUrl = url + "&cursor=" + response.NextCursorKey
	}
}
