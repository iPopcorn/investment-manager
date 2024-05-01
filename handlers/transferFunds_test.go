package handlers_test

import (
	"testing"

	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/iPopcorn/investment-manager/infrastructure"
	testutils "github.com/iPopcorn/investment-manager/test-utils"
)

func TestTransferFunds(t *testing.T) {
	t.Run("Fails to transfer if number of args unexpected", func(t *testing.T) {
		testHttpClient := &testutils.TestHttpClient{
			GetResponse: "",
		}
		testInternalClient := infrastructure.InvestmentManagerInternalHttpClientFactory(testHttpClient, "")
		testHandler := handlers.TransferFundsHandlerFactory(testInternalClient)

		args := []string{}
		err := testHandler(testutils.TestCmd, args)

		assertError(err, args, t)

		args = []string{"test1"}
		err = testHandler(testutils.TestCmd, args)

		assertError(err, args, t)

		args = []string{"test1", "test2"}
		err = testHandler(testutils.TestCmd, args)

		assertError(err, args, t)

		args = []string{"test1", "test2", "test3", "test4"}
		err = testHandler(testutils.TestCmd, args)

		assertError(err, args, t)
	})

	t.Run("Fails to transfer if amount can't be converted to float", func(t *testing.T) {
		testHttpClient := &testutils.TestHttpClient{
			GetResponse: "",
		}
		testInternalClient := infrastructure.InvestmentManagerInternalHttpClientFactory(testHttpClient, "")
		testHandler := handlers.TransferFundsHandlerFactory(testInternalClient)

		args := []string{"sender", "receiver", "twenty"}
		err := testHandler(testutils.TestCmd, args)

		assertError(err, args, t)
	})
}

func assertError(err error, args []string, t *testing.T) {
	t.Helper()

	if err == nil {
		t.Fatalf("Expected error but did not receive one.\nHandler called with %v\n", args)
	}
}
