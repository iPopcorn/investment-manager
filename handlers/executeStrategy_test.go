package handlers_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/iPopcorn/investment-manager/infrastructure"
	testutils "github.com/iPopcorn/investment-manager/test-utils"
)

type TestExecuteStrategyHttpClient struct {
	counter int
}

type TestReader struct {
	data []byte
}

func (testClient *TestExecuteStrategyHttpClient) Do(req *http.Request) (*http.Response, error) {
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
		testInternalClient := infrastructure.InvestmentManagerInternalHttpClientFactory(testHttpClient, "")
		testHandler := handlers.ExecuteStrategyHandlerFactory(testInternalClient)
		args := []string{"test", "hodl", "eth"}

		err := testHandler(testutils.TestCmd, args)

		if err != nil {
			t.Fatalf("Received error but did not expect one\n%v", err)
		}

		err = testHandler(testutils.TestCmd, args)

		if err == nil {
			t.Fatalf("Expected error but did not receive one")
		}
	})
}
