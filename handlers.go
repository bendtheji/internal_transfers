package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type CreateAccountRequest struct {
	ID      int     `json:"account_id"`
	Balance float64 `json:"initial_balance"`
}

type GetAccountResponse struct {
	ID      int     `json:"account_id"`
	Balance float64 `json:"balance"`
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	db, err := ConnectToDB(*dbConfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var req CreateAccountRequest
	json.NewDecoder(r.Body).Decode(&req)

	err = CreateAccount(db, req.ID, req.Balance)
	if err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Account created")
}

func getAccountHandler(w http.ResponseWriter, r *http.Request) {
	db, err := ConnectToDB(*dbConfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Get the 'id' parameter from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert 'id' to an integer
	accountId, err := strconv.Atoi(idStr)

	// Call the GetUser function to fetch the user data from the database
	account, err := GetAccount(db, accountId)
	if err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	accountResponse := GetAccountResponse{
		account.id,
		account.balance,
	}

	// Convert the user object to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accountResponse)
}

func createTransactionHandler(w http.ResponseWriter, r *http.Request) {}
