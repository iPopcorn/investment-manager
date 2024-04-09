package cmd_test

import (
	"testing"

	testutils "github.com/iPopcorn/investment-manager/test-utils"
)

func TestPortfolioDetails(t *testing.T) {
	t.Run("Returns an error when portfolio-details command called with no args", func(t *testing.T) {
		args := []string{"portfolio-details"}
		expectError(testutils.ExecuteCommand(args), t)
	})

	t.Run("Returns an error when portfolio-details command called with too many args", func(t *testing.T) {
		args := []string{"portfolio-details", "default", "crypto"}

		expectError(testutils.ExecuteCommand(args), t)
	})

	t.Run("Returns an error when portfolio-details command cannot find given portfolio", func(t *testing.T) {
		args := []string{"portfolio-details", "portfolio-doesnt-exist"}

		expectError(testutils.ExecuteCommand(args), t)
	})
}

func expectError(err error, t *testing.T) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected error but received none")
	}
}
