package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
	"github.com/spf13/cobra"
)

func ExecuteStrategyHandlerFactory(client *infrastructure.InvestmentManagerInternalHttpClient) CobraCommandHandler {
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
	if strings.ToUpper(strategy) != string(types.HODL) {
		return fmt.Errorf("Invalid strategy\nGiven: %q Expected: %q\n", strategy, string(types.HODL))
	}

	if strings.ToUpper(currency) != string(types.ETH) {
		return fmt.Errorf("Invalid currency\nGiven: %q Expected: %q\n", currency, string(types.ETH))
	}

	request := &types.ExecuteStrategyRequest{
		Portfolio: portfolio,
		Strategy:  types.HODL,
		Currency:  types.ETH,
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
