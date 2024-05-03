package cmd

import (
	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/spf13/cobra"
)

var transferFundsCommand = &cobra.Command{
	Use:   "transfer-funds sender receiver amount",
	Short: "Transfer funds from the sender to the receiver",
	Long: `Transfer funds from the sender portfolio to the receiver portfolio. 
An error is thrown if incorrect arguments given. 
Refer to the portfolios by name.
Names are case sensitive.
Use 'portfolio' command to see list of portfolios.
example: 'transfer-funds default test 10'`,
	RunE: nil,
}

func init() {
	internalHttpClient := infrastructure.GetDefaultInvestmentManagerInternalHttpClient()
	transferFundsHandler := handlers.TransferFundsHandlerFactory(internalHttpClient)
	transferFundsCommand.RunE = transferFundsHandler

	rootCmd.AddCommand(transferFundsCommand)
}
