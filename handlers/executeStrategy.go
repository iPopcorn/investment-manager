package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/spf13/cobra"
)

func ExecuteStrategyHandlerFactory(client *infrastructure.InvestmentManagerInternalHttpClient) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Printf("ExecuteStrategy called\nargs: %v", args)

		if len(args) != 3 {
			return fmt.Errorf("Expected 3 args, received %d args", len(args))
		}

		portfolio := args[0]
		strategy := args[1]
		currency := args[2]

		return executeStrategy(portfolio, strategy, currency, client)
	}
}

func executeStrategy(portfolio, strategy, currency string, client *infrastructure.InvestmentManagerInternalHttpClient) error {
	request := &types.ExecuteStrategyRequest{
		Portfolio: portfolio,
		Strategy:  strategy,
		Currency:  currency,
	}
	serializedRequest, err := json.Marshal(request)

	if err != nil {
		fmt.Printf("Failed to serialize request: %v", err)
		return err
	}

	_, err = client.Post("/execute-strategy", serializedRequest)

	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return err
	}

	return nil
}
