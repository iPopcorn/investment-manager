package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/spf13/cobra"
)

func HandlePortfolioFactory(client *infrastructure.InvestmentManagerInternalHttpClient) CobraCommandHandler {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Printf("portfolio called\nargs: %v\n", args)
		fmt.Printf("listing portfolios...\n")
		portfolios, err := listPortfolios(client)

		if err != nil {
			return err
		}

		displayPortfolios(portfolios)
		return nil
	}
}

func displayPortfolios(portfolioResponse *types.PortfolioResponse) {
	fmt.Println("Portfolios:")
	for i, p := range portfolioResponse.Portfolios {
		fmt.Printf("%d)\n", i+1)
		fmt.Printf(" Name: %s\n", p.Name)
		fmt.Printf(" Type: %s\n", p.Type)
		fmt.Printf(" UUID: %s\n", p.Uuid)
		fmt.Printf(" Is Deleted?: %t\n", p.Deleted)
	}
}

func listPortfolios(client *infrastructure.InvestmentManagerInternalHttpClient) (*types.PortfolioResponse, error) {
	path := "/portfolios"

	httpResponse, err := client.Get(path)

	if err != nil {
		return nil, fmt.Errorf("Error getting portfolios from api: \n%v\n", err)
	}

	var resp types.PortfolioResponse
	json.Unmarshal(httpResponse, &resp)

	return &resp, nil
}
