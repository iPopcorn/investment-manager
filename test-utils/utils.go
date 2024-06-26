package testutils

import (
	"bytes"
	"io"
	"net/http"

	"github.com/iPopcorn/investment-manager/cmd"
	"github.com/spf13/cobra"
)

var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "test command",
	Long:  `dummy command to use for testing`,
	RunE:  nil,
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

type TestHttpClient struct {
	GetResponse string
}

type TestReader struct {
	data []byte
}

func (testClient TestHttpClient) Do(req *http.Request) (*http.Response, error) {
	resp := http.Response{
		Body: io.NopCloser(bytes.NewBufferString(testClient.GetResponse)),
	}
	return &resp, nil
}
