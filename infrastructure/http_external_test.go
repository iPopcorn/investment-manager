package infrastructure_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/iPopcorn/investment-manager/infrastructure"
)

type testHttpClient struct {
	getResponse string
}

type testReader struct {
	data []byte
}

func (testClient testHttpClient) Do(req *http.Request) (*http.Response, error) {
	resp := http.Response{
		Body: io.NopCloser(bytes.NewBufferString(testClient.getResponse)),
	}
	return &resp, nil
}

func TestGet(t *testing.T) {
	t.Run("Gets data from the given URL", func(t *testing.T) {
		expected := "Called Do()"
		httpClient := testHttpClient{getResponse: expected}
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
		httpClient := testHttpClient{}
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
