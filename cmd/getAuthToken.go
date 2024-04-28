package cmd

import (
	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/spf13/cobra"
)

var getAuthTokenCmd = &cobra.Command{
	Use:   "get-token",
	Short: "get an auth token",
	Long: `Get an auth token. 
Use to manually make requests to coinbase.`,
	RunE: handlers.HandleGetAuthToken,
}

func init() {
	rootCmd.AddCommand(getAuthTokenCmd)
}
