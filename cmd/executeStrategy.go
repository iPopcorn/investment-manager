package cmd

import (
	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/spf13/cobra"
)

var executeStrategyCmd = &cobra.Command{
	Use:   "execute-strategy portfolio strategy currency",
	Short: "Execute a specified trading strategy against a given portfolio",
	Long: `Execute a trading strategy against a given portfolio. 
If no portfolio is given, an error is thrown. 
Use 'portfolio' command to see list of portfolios.
Supported strategies:
HODL
Supported currencies:
ETH
example: 'execute-strategy test hodl eth'`,
	RunE: nil,
}

func init() {
	internalHttpClient := infrastructure.GetDefaultInvestmentManagerInternalHttpClient()

	executeStrategyHandler := handlers.ExecuteStrategyHandlerFactory(internalHttpClient)
	executeStrategyCmd.RunE = executeStrategyHandler

	rootCmd.AddCommand(executeStrategyCmd)
}
