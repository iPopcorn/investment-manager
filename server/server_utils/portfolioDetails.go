package server_utils

import (
	"log"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/iPopcorn/investment-manager/types/mappers"
)

func PortfolioDetails(client *infrastructure.InvestmentManagerExternalHttpClient, portfolioID string) (*types.PortfolioDetailsResponse, error) {
	url := "https://api.coinbase.com/api/v3/brokerage/portfolios"
	url = url + "/" + portfolioID
	resp, err := client.Get(url)

	if err != nil {
		log.Printf("Error retrieving portfolio details from URL: %q\nError: %v", url, err)
		return nil, err
	}

	return mappers.MapPortfolioDetailsResponse(resp)
}
