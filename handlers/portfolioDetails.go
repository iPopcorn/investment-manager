package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/iPopcorn/investment-manager/util"
	"github.com/spf13/cobra"
)

func PortfolioDetailsFactory(internalHttpClient *infrastructure.InvestmentManagerInternalHttpClient) CobraCommandHandler {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Expected an argument but did not receive one")
		}

		if len(args) > 1 {
			return fmt.Errorf("Expected 1 arg, received %d args", len(args))
		}

		details, err := getPortfolioDetails(args[0], internalHttpClient)

		if err != nil {
			return err
		}

		showPortfolioDetails(details)
		return err
	}
}

func showPortfolioDetails(details *types.PortfolioDetailsResponse) {
	name := details.Breakdown.Portfolio.Name
	cashBalance := details.Breakdown.PortfolioBalances.TotalCashEquivalentBalance
	totalBalance := details.Breakdown.PortfolioBalances.TotalBalance

	fmt.Println("Portfolio Details")
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Total Value: %s %s\n", totalBalance.Value, totalBalance.Currency)
	fmt.Printf("Amount available for trade: %s %s\n", cashBalance.Value, cashBalance.Currency)
}

func getPortfolioDetails(portfolioName string, client *infrastructure.InvestmentManagerInternalHttpClient) (*types.PortfolioDetailsResponse, error) {
	portfolios, err := listPortfolios(client)

	if err != nil {
		return nil, err
	}

	foundPortfolio, err := findPortfolio(portfolioName, portfolios.Portfolios)

	if err != nil {
		return nil, err
	}

	portfolioID := foundPortfolio.Uuid
	path := "/portfolios/" + portfolioID

	httpResponse, err := client.Get(path)

	if err != nil {
		return nil, fmt.Errorf("Error getting portfolio details from api: \n%v\n", err)
	}

	err = util.HandleErrorResponse(httpResponse)

	if err != nil {
		fmt.Println("Failed to retrieve portfolio details")
		return nil, err
	}

	var portfolioDetailsResponse types.PortfolioDetailsResponse
	err = json.Unmarshal(httpResponse, &portfolioDetailsResponse)

	if err != nil {
		fmt.Println("failed to parse response")
		return nil, err
	}

	return &portfolioDetailsResponse, nil
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
