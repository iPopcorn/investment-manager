package cmd_test

import (
	"bytes"
	"testing"

	"github.com/iPopcorn/investment-manager/cmd"
)

func TestPortfolioDetails(t *testing.T) {
	t.Run("Returns an error when portfolio-details command called with no args", func(t *testing.T) {
		// capture output in program by using a buffer
		actual := new(bytes.Buffer)

		// initialize root command for test context
		testCmd := cmd.NewRootCmd()

		// set output of root command to buffer
		testCmd.SetOut(actual)
		testCmd.SetErr(actual)

		// pass args to root command - this is how you call sub commands
		testCmd.SetArgs([]string{"portfolio-details"})

		// execute command and expect error
		err := testCmd.Execute()

		if err == nil {
			t.Fatalf("Expected error but received none")
		}
	})
}
