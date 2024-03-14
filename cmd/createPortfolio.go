package cmd

import (
	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/spf13/cobra"
)

var createPortfolioCmd = &cobra.Command{
	Use:   "create-portfolio name",
	Short: "Create a new portfolio",
	Long: `Create a new portfolio. 
If no name is given, an error is thrown.`,
	RunE: handlers.CreatePortfolio,
}

func init() {
	rootCmd.AddCommand(createPortfolioCmd)
}
