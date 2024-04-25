package api

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"

	dbPackage "github.com/bendtheji/internal_transfers/db"
	apiError "github.com/bendtheji/internal_transfers/errors"
)

type CreateTransactionRequest struct {
	SourceAccountID      int    `json:"source_account_id"`
	DestinationAccountID int    `json:"destination_account_id"`
	TransactionID        string `json:"transaction_id"`
	Amount               string `json:"amount"`
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbPackage.ConnectToDB(*dbPackage.DbConfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var req CreateTransactionRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apiError.HandleApiError(w, apiError.HandleError(fmt.Errorf("%w: %w", apiError.ReqUnmarshalTypeErr, err)))
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		apiError.HandleApiError(w, apiError.HandleError(fmt.Errorf("%w: %v", apiError.InvalidTransactionAmountErr, req.Amount)))
		return
	}

	amount = amount.Truncate(2)
	if amount.LessThanOrEqual(decimal.Zero) {
		apiError.HandleApiError(w, apiError.HandleError(fmt.Errorf("%w: %v", apiError.InvalidTransactionAmountErr, req.Amount)))
		return
	}

	transaction := dbPackage.Transaction{
		TransactionID:        req.TransactionID,
		SourceAccountID:      req.SourceAccountID,
		DestinationAccountID: req.DestinationAccountID,
		Amount:               amount,
	}

	err = dbPackage.CreateTransaction(r.Context(), db, &transaction)
	if err != nil {
		apiError.HandleApiError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Transaction completed successfully")
}
