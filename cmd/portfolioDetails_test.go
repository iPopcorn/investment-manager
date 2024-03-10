package cmd_test

import (
	"bytes"
	"testing"

	"github.com/iPopcorn/investment-manager/cmd"
)

func TestPortfolioDetails(t *testing.T) {
	t.Run("Returns an error when portfolio-details command called with no args", func(t *testing.T) {
		args := []string{"portfolio-details"}
		expectError(ExecuteCommand(args), t)
	})

	t.Run("Returns an error when portfolio-details command called with too many args", func(t *testing.T) {
		args := []string{"portfolio-details", "default", "crypto"}

		expectError(ExecuteCommand(args), t)
	})
}

func ExecuteCommand(args []string) error {
	// capture output in program by using a buffer
	output := new(bytes.Buffer)

	// initialize root command for test context
	testCmd := cmd.NewRootCmd()

	// set output of root command to buffer
	testCmd.SetOut(output)
	testCmd.SetErr(output)

	// pass args to root command - this is how you call sub commands
	testCmd.SetArgs(args)

	// execute command
	return testCmd.Execute()
}

func expectError(err error, t *testing.T) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected error but received none")
	}
}
