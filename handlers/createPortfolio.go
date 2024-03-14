package handlers

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CreatePortfolio(cmd *cobra.Command, args []string) error {
	fmt.Printf("CreatePortfolio called\nargs: %v", args)

	if len(args) != 1 {
		return fmt.Errorf("Expected 1 arg, received %d args", len(args))
	}

	return nil
}
