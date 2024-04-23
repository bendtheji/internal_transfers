package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

var (
	dbConfig *EnvDBConfig
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConfig = NewEnvDBConfig()
	r := mux.NewRouter()

	r.HandleFunc("/accounts", createAccountHandler).Methods("POST")
	r.HandleFunc("/accounts/{id}", getAccountHandler).Methods("GET")
	r.HandleFunc("/transactions", createTransactionHandler).Methods("POST")

	log.Println("Server listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
