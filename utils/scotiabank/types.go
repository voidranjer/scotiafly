package scotiabank

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	NextCursorKey string       `json:"nextCursorKey"`
	Transactions  Transactions `json:"data"`
}

type Transaction struct {
	CleanDescription  string            `json:"cleanDescription"`
	TransactionDate   string            `json:"transactionDate"`
	TransactionType   string            `json:"transactionType"` // "DEBIT" or "CREDIT"
	TransactionAmount TransactionAmount `json:"transactionAmount"`
	Category          Category          `json:"category"`
	Id                string            `json:"id"`
}

type TransactionAmount struct {
	Amount       float32 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
}

type Category struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

/* Transactions */

type Transactions []Transaction

func (t *Transactions) UnmarshalJSON(data []byte) error {
	// Try first as []Transaction
	var direct []Transaction
	if err := json.Unmarshal(data, &direct); err == nil {
		*t = direct
		return nil
	}

	// Try as object with "settled" field
	var nested struct {
		Settled []Transaction `json:"settled"`
	}
	if err := json.Unmarshal(data, &nested); err == nil {
		*t = nested.Settled
		return nil
	}

	return fmt.Errorf("unsupported JSON structure for Transactions")
}
