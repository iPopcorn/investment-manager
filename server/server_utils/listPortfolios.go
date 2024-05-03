package server_utils

import (
	"fmt"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/iPopcorn/investment-manager/types/mappers"
)

func ListPortfolios(client *infrastructure.InvestmentManagerExternalHttpClient) (*types.PortfolioResponse, error) {
	url := "https://api.coinbase.com/api/v3/brokerage/portfolios"

	resp, err := client.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving portfolios from URL: %q\nError: %v", url, err)
	}

	return mappers.MapPortfolioResponse(resp)
}
