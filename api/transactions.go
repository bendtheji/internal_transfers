package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	dbPackage "github.com/bendtheji/internal_transfers/db"
	apiError "github.com/bendtheji/internal_transfers/errors"
)

type CreateTransactionRequest struct {
	SourceAccountID      int     `json:"source_account_id"`
	DestinationAccountID int     `json:"destination_account_id"`
	TransactionID        string  `json:"transaction_id"`
	Amount               float64 `json:"amount"`
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbPackage.ConnectToDB(*dbPackage.DbConfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var req CreateTransactionRequest
	json.NewDecoder(r.Body).Decode(&req)

	transaction := dbPackage.Transaction{
		TransactionID:        req.TransactionID,
		SourceAccountID:      req.SourceAccountID,
		DestinationAccountID: req.DestinationAccountID,
		Amount:               req.Amount,
	}

	err = dbPackage.CreateTransaction(r.Context(), db, &transaction)
	if err != nil {
		apiError.HandleApiError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Transaction completed successfully")
}
