package infrastructure_test

import (
	"fmt"
	"testing"

	"github.com/iPopcorn/investment-manager/infrastructure"
	testutils "github.com/iPopcorn/investment-manager/test-utils"
)

func TestGet(t *testing.T) {
	t.Run("Gets data from the given URL", func(t *testing.T) {
		expected := "Called Do()"
		httpClient := testutils.TestHttpClient{GetResponse: expected}
		testClient := infrastructure.InvestmentManagerExternalHttpClient{
			HttpClient: httpClient,
		}
		data, err := testClient.Get("https://api.coinbase.com/api/v3/brokerage/portfolios")

		if err != nil {
			t.Fatalf("Unexpected error\n%v", err)
		}

		actual := string(data)

		if actual != expected {
			t.Fatalf("Expected: '%s'\nActual: '%s'\n", expected, actual)
		}
	})

	t.Run("Returns error if url is invalid", func(t *testing.T) {
		httpClient := testutils.TestHttpClient{}
		testClient := infrastructure.InvestmentManagerExternalHttpClient{
			HttpClient: httpClient,
		}
		data, err := testClient.Get("Invalid URL")

		if err == nil {
			t.Fatalf("expected error\n")
		}

		fmt.Printf("data: %s\n", string(data))
	})
}
