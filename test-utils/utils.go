package testutils

import (
	"bytes"

	"github.com/iPopcorn/investment-manager/cmd"
)

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
