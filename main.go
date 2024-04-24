package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/bendtheji/internal_transfers/api"
	"github.com/bendtheji/internal_transfers/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDbConfig()
	r := mux.NewRouter()

	r.HandleFunc("/accounts", api.CreateAccountHandler).Methods("POST")
	r.HandleFunc("/accounts/{id}", api.GetAccountHandler).Methods("GET")
	r.HandleFunc("/transactions", api.CreateTransactionHandler).Methods("POST")

	log.Println("Server listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
