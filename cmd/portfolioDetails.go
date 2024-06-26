package cmd

import (
	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/spf13/cobra"
)

var portfolioDetailsCmd = &cobra.Command{
	Use:   "portfolio-details",
	Short: "Display details about a given portfolio",
	Long: `Display details about a given portfolio. 
If no portfolio is given, display details of the first
portfolio in the list. use 'portfolio' command to see list of portfolios.`,
	RunE: nil,
}

func init() {
	client := infrastructure.GetDefaultInvestmentManagerInternalHttpClient()
	handler := handlers.PortfolioDetailsFactory(client)
	portfolioDetailsCmd.RunE = handler
	rootCmd.AddCommand(portfolioDetailsCmd)
}
