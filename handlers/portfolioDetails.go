package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/spf13/cobra"
)

func PortfolioDetails(cmd *cobra.Command, args []string) error {
	fmt.Printf("Called portfolio-details\nargs: %v\n", args)

	if len(args) == 0 {
		return errors.New("Expected an argument but did not receive one")
	}

	if len(args) > 1 {
		return fmt.Errorf("Expected 1 arg, received %d args", len(args))
	}

	err := listPortfolioDetails(args[0])
	return err
}

func listPortfolioDetails(portfolioName string) error {
	portfolios, err := listPortfolios()

	if err != nil {
		return err
	}

	foundPortfolio, err := findPortfolio(portfolioName, portfolios.Portfolios)

	if err != nil {
		return err
	}

	portfolioID := foundPortfolio.Uuid

	url := "https://api.coinbase.com/api/v3/brokerage/portfolios/" + portfolioID

	httpClient := infrastructure.InvestmentManagerHTTPClient{
		HttpClient: &http.Client{},
	}

	httpResponse, err := httpClient.Get(url)

	if err != nil {
		return fmt.Errorf("Error getting portfolios from api: \n%v\n", err)
	}

	var portfolioDetailsResponse types.PortfolioDetailsResponse
	err = json.Unmarshal(httpResponse, &portfolioDetailsResponse)

	fmt.Printf("%+v", portfolioDetailsResponse)
	return nil
}

func findPortfolio(name string, portfolios []types.Portfolio) (*types.Portfolio, error) {
	var foundPortfolio *types.Portfolio = nil
	for _, portfolio := range portfolios {
		if strings.ToLower(name) == strings.ToLower(portfolio.Name) {
			foundPortfolio = &portfolio
			break
		}
	}

	if foundPortfolio == nil {
		return nil, fmt.Errorf("Could not find portfolio with the name: %s", name)
	}

	return foundPortfolio, nil
}
