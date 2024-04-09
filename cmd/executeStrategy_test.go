package cmd_test

import (
	"testing"

	testutils "github.com/iPopcorn/investment-manager/test-utils"
)

func TestExecuteStrategy(t *testing.T) {
	t.Run("Fails to execute a strategy that is already running", func(t *testing.T) {
		args := []string{"execute-strategy", "test", "hodl", "eth"}
		err := testutils.ExecuteCommand(args)

		if err != nil {
			t.Fatalf("Received error but did not expect one\n%v", err)
		}

		err = testutils.ExecuteCommand(args)

		if err == nil {
			t.Fatalf("Expected error but did not receive one")
		}
	})
}
