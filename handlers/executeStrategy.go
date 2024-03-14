package handlers

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ExecuteStrategy(cmd *cobra.Command, args []string) error {
	fmt.Printf("ExecuteStrategy called\nargs: %v", args)

	if len(args) != 3 {
		return fmt.Errorf("Expected 3 args, received %d args", len(args))
	}

	return nil
}
