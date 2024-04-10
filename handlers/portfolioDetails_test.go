package handlers_test

import (
	"testing"

	"github.com/iPopcorn/investment-manager/handlers"
	"github.com/iPopcorn/investment-manager/infrastructure"
	testutils "github.com/iPopcorn/investment-manager/test-utils"
)

func TestPortfolioDetails(t *testing.T) {
	t.Run("Returns an error when portfolio-details command called with no args", func(t *testing.T) {
		args := []string{"portfolio-details"}
		testHandler := getTestPortfolioDetailsHandler()

		expectError(testHandler(testutils.TestCmd, args), t)
	})

	t.Run("Returns an error when portfolio-details command called with too many args", func(t *testing.T) {
		args := []string{"portfolio-details", "default", "crypto"}
		testHandler := getTestPortfolioDetailsHandler()

		expectError(testHandler(testutils.TestCmd, args), t)
	})

	t.Run("Returns an error when portfolio-details command cannot find given portfolio", func(t *testing.T) {
		args := []string{"portfolio-details", "portfolio-doesnt-exist"}
		testHandler := getTestPortfolioDetailsHandler()

		expectError(testHandler(testutils.TestCmd, args), t)
	})
}

func getTestPortfolioDetailsHandler() handlers.CobraCommandHandler {
	testHttpClient := &testutils.TestHttpClient{GetResponse: ""}
	testInternalClient := infrastructure.GetInjectedInvestmentManagerInternalHttpClient(testHttpClient, "test")

	return handlers.PortfolioDetailsFactory(testInternalClient)
}

func expectError(err error, t *testing.T) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected error but received none")
	}
}
