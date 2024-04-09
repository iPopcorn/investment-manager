package handlers_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test command",
	Long:  `dummy command to use for testing`,
	RunE:  nil,
}

type TestExecuteStrategyHttpClient struct {
	counter int
}

type TestReader struct {
	data []byte
}

func (testClient *TestExecuteStrategyHttpClient) Do(req *http.Request) (*http.Response, error) {
	fmt.Printf("Do()\ncounter: %d\n", testClient.counter)
	testClient.counter = testClient.counter + 1
	resp := http.Response{
		Body: io.NopCloser(bytes.NewBufferString("")),
	}

	if testClient.counter > 1 {
		return nil, fmt.Errorf("strategy is being executed\ncounter: %d", testClient.counter)
	}

	return &resp, nil
}

func TestExecuteStrategy(t *testing.T) {
	t.Run("Fails to execute a strategy that is already running", func(t *testing.T) {

		testHttpClient := &TestExecuteStrategyHttpClient{counter: 0}
		testInternalClient := infrastructure.GetInjectedInvestmentManagerInternalHttpClient(testHttpClient, "")
		testHandler := handlers.ExecuteStrategyHandlerFactory(testInternalClient)
		args := []string{"test", "hodl", "eth"}

		err := testHandler(testCmd, args)

		if err != nil {
			t.Fatalf("Received error but did not expect one\n%v", err)
		}

		err = testHandler(testCmd, args)

		if err == nil {
			t.Fatalf("Expected error but did not receive one")
		}
	})
}
