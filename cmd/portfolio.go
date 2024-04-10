/*
Copyright © 2024 Taylor Petrillo <taypetrillo@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/iPopcorn/investment-manager/infrastructure"
)

// portfolioCmd represents the portfolio command
var portfolioCmd = &cobra.Command{
	Use:   "portfolio",
	Short: "display information about your portfolio(s)",
	Long: `Display information about your portfolio(s)
Calling 'portfolio' will list all portfolios associated with your account`,
	RunE: nil,
}

func init() {
	client := infrastructure.GetDefaultInvestmentManagerInternalHttpClient()
	handler := handlers.HandlePortfolioFactory(client)
	portfolioCmd.RunE = handler

	rootCmd.AddCommand(portfolioCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// portfolioCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	portfolioCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
