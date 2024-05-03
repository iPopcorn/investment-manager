package main

import (
	"log"
	"net/http"

	"github.com/iPopcorn/investment-manager/server"
)

func main() {
	address := "127.0.0.1:5000"
	server := server.GetDefaultInvestmentManagerHTTPServer()

	log.Printf("Listening at %s\n", address)
	log.Fatal(http.ListenAndServe(address, server))
}
