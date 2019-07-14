// Golang API Boilerplate
// Source: https://github.com/orgmatileg/golang-api-boilerplate-crud
package main

import (
	"fmt"
	"golang-api-boilerplate-crud/database"
	"golang-api-boilerplate-crud/routes"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const appVer = "1.0.0.071419"

func init() {
	// If Env DEPLOY is not "PRODUCTION"
	// Load development .env
	if e := os.Getenv("APP_DEPLOY"); e != "PRODUCTION" {
		loadDotEnvDevelopment()
	}
}

func main() {
	log.Printf("Application running on version: %s\n", appVer)
	database.Initiation()
	serveHTTP()
}

func serveHTTP() {
	r := routes.Initiation()
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST_ADDRESS"), os.Getenv("PORT"))
	log.Println("App running on", addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("could not serve http server on port 8080: %s", err.Error())
	}
}

func loadDotEnvDevelopment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err.Error())
	}
}
