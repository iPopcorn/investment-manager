package main

import (
	"log"
	"net/http"

	"github.com/iPopcorn/investment-manager/server"
)

func main() {
	server := server.GetInvestmentManagerHTTPServer()
	log.Fatal(http.ListenAndServe("127.0.0.1:5000", server))
}
