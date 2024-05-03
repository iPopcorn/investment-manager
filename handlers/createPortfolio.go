package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/iPopcorn/investment-manager/util"
	"github.com/spf13/cobra"
)

func CreatePortfolio(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected 1 arg, received %d args", len(args))
	}

	newPortfolio, err := createPortfolio(args[0])

	if err != nil {
		return err
	}

	displayCreatedPortfolio(newPortfolio)
	return nil
}

func displayCreatedPortfolio(newPortfolio *types.PortfolioCreatedResponse) {
	p := newPortfolio.Portfolio
	fmt.Println("Created a new portfolio:")
	fmt.Printf(" Name: %s\n", p.Name)
	fmt.Printf(" Type: %s\n", p.Type)
	fmt.Printf(" UUID: %s\n", p.Uuid)
	fmt.Printf(" Is Deleted?: %t\n", p.Deleted)
}

func createPortfolio(name string) (*types.PortfolioCreatedResponse, error) {
	path := "/portfolios"
	internalClient := infrastructure.GetDefaultInvestmentManagerInternalHttpClient()

	request := []byte(fmt.Sprintf(`{
		"name": "%s"
	}`, name))

	httpResponse, err := internalClient.Post(path, request)

	if err != nil {
		return nil, fmt.Errorf("Error creating portfolio: \n%v\n", err)
	}

	err = util.HandleErrorResponse(httpResponse)

	if err != nil {
		fmt.Println("Failed to create portfolio")
		return nil, err
	}

	var resp types.PortfolioCreatedResponse
	err = json.Unmarshal(httpResponse, &resp)

	if err != nil {
		fmt.Println("Failed to parse response")
		return nil, err
	}

	return &resp, nil
}
