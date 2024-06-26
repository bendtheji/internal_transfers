package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"

	dbPackage "github.com/bendtheji/internal_transfers/db"
	apiError "github.com/bendtheji/internal_transfers/errors"
)

type CreateAccountRequest struct {
	ID      int    `json:"account_id"`
	Balance string `json:"initial_balance"`
}

type GetAccountResponse struct {
	ID      int    `json:"account_id"`
	Balance string `json:"balance"`
}

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbPackage.ConnectToDB(*dbPackage.DbConfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var req CreateAccountRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		apiError.HandleApiError(w, apiError.HandleError(fmt.Errorf("%w: %w", apiError.ReqUnmarshalTypeErr, err)))
		return
	}

	// check that ID and balance is valid values
	if req.ID <= 0 {
		apiError.HandleApiError(w, apiError.HandleError(fmt.Errorf("%w: %v", apiError.InvalidAccountIDErr, req.ID)))
		return
	}

	balance, err := decimal.NewFromString(req.Balance)
	if err != nil {
		apiError.HandleApiError(w, apiError.HandleError(fmt.Errorf("%w: %v", apiError.InvalidInitialBalanceErr, req.Balance)))
		return
	}

	balance = balance.Truncate(2)
	if balance.LessThanOrEqual(decimal.Zero) {
		apiError.HandleApiError(w, apiError.HandleError(fmt.Errorf("%w: %v", apiError.InvalidInitialBalanceErr, req.Balance)))
		return
	}

	err = dbPackage.CreateAccount(r.Context(), db, req.ID, balance)
	if err != nil {
		apiError.HandleApiError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Account created")
}

func GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbPackage.ConnectToDB(*dbPackage.DbConfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Get the 'id' parameter from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert 'id' to an integer
	accountId, err := strconv.Atoi(idStr)
	if err != nil {
		apiError.HandleApiError(w, apiError.HandleError(fmt.Errorf("%w: %v", apiError.InvalidAccountIDErr, idStr)))
		return
	}

	// Call the GetUser function to fetch the user data from the database
	account, err := dbPackage.GetAccount(r.Context(), db, accountId)
	if err != nil {
		apiError.HandleApiError(w, err)
		return
	}

	accountResponse := GetAccountResponse{
		account.ID,
		account.Balance.String(),
	}

	// Convert the user object to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accountResponse)
}
