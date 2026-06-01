package main

import (
	"log"
	"net/http"
	"os"

	"cnzamnt/backend/internal/api"
	"cnzamnt/backend/internal/db"
)

func main() {
	addr := getenv("CNZAMNT_ADDR", ":8080")
	dbPath := getenv("CNZAMNT_DB_PATH", "data/cnzamnt.db")

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	server := api.New(database)
	log.Printf("CnzAMnt API listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
