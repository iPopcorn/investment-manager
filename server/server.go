package server

import (
	"log"
	"net/http"

	"github.com/iPopcorn/investment-manager/infrastructure"
)

type InvestmentManagerHTTPServer struct {
	client infrastructure.InvestmentManagerExternalHttpClient
}

func (s *InvestmentManagerHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %v\n", r)
	w.Header().Set("Content-Type", "application/json")
	url := "https://api.coinbase.com/api/v3/brokerage/portfolios"
	resp, err := s.client.Get(url)
	if err != nil {
		log.Printf("Error retrieving portfolios from URL: %q\nError: %v", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)

	if err != nil {
		log.Println("Failed to write response to writer")
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("Request handled successfully!")

	return
}

func GetInvestmentManagerHTTPServer() *InvestmentManagerHTTPServer {
	httpClient := infrastructure.InvestmentManagerExternalHttpClient{
		HttpClient: &http.Client{},
	}

	return &InvestmentManagerHTTPServer{
		client: httpClient,
	}
}
