package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/spf13/cobra"
)

type PortfolioResponse struct {
	Portfolios []Portfolio `json:"portfolios"`
}

type Portfolio struct {
	Name    string `json:"name"`
	Uuid    string `json:"uuid"`
	Type    string `json:"type"`
	Deleted bool   `json:"deleted"`
}

func HandlePortfolio(cmd *cobra.Command, args []string) {
	fmt.Printf("portfolio called\nargs: %v\n", args)
	fmt.Printf("listing portfolios...\n")
	listPortfolios()
}

func listPortfolios() {
	url := "https://api.coinbase.com/api/v3/brokerage/portfolios"

	httpClient := infrastructure.InvestmentManagerHTTPClient{
		HttpClient: &http.Client{},
	}

	httpResponse, err := httpClient.Get(url)

	if err != nil {
		fmt.Errorf("Error getting portfolios from api: \n%v\n", err)
	}

	var resp PortfolioResponse
	json.Unmarshal(httpResponse, &resp)

	fmt.Println("Portfolios:")
	for i, p := range resp.Portfolios {
		fmt.Printf("%d)\n", i+1)
		fmt.Printf(" Name: %s\n", p.Name)
		fmt.Printf(" Type: %s\n", p.Type)
		fmt.Printf(" UUID: %s\n", p.Uuid)
		fmt.Printf(" Is Deleted?: %t\n", p.Deleted)
	}
}
