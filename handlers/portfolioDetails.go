package handlers

import (
	"fmt"

	"github.com/spf13/cobra"
)

func PortfolioDetails(cmd *cobra.Command, args []string) {
	fmt.Printf("Called portfolio-details\nargs: %v\n", args)
}

func listPortfolioDetails() {

}
