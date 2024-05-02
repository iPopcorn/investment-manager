package cmd

import (
	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/spf13/cobra"
)

var transferFundsCommand = &cobra.Command{
	Use:   "transfer-funds sender receiver amount",
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
	transferFundsHandler := handlers.TransferFundsHandlerFactory(internalHttpClient)
	transferFundsCommand.RunE = transferFundsHandler

	rootCmd.AddCommand(transferFundsCommand)
}
