package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
)

type testHttpClient struct {
	getResponse []byte
}

type testReader struct {
	data []byte
}

func (testClient testHttpClient) Do(req *http.Request) (*http.Response, error) {
	resp := http.Response{
		Body: io.NopCloser(bytes.NewReader(testClient.getResponse)),
	}
	return &resp, nil
}

func getTestServer(expectedResponse []byte) *InvestmentManagerHTTPServer {
	testHttpClient := testHttpClient{
		getResponse: expectedResponse,
	}
	testInvestmentManagerHTTPClient := infrastructure.InvestmentManagerHTTPClient{
		HttpClient: testHttpClient,
	}
	return &InvestmentManagerHTTPServer{
		client: testInvestmentManagerHTTPClient,
	}
}

func TestGETPortfolios(t *testing.T) {
	t.Run("Gets user's portfolios from coinbase", func(t *testing.T) {
		portfolio1 := types.Portfolio{
			Name:    "Test One",
			Uuid:    "test-portfolio-1",
			Type:    "type1",
			Deleted: false,
		}
		portfolio2 := types.Portfolio{
			Name:    "Test Two",
			Uuid:    "test-portfolio-2",
			Type:    "type2",
			Deleted: false,
		}
		expectedResponse := &types.PortfolioResponse{
			Portfolios: []types.Portfolio{portfolio1, portfolio2},
		}

		serializedExpectedResponse, err := json.Marshal(expectedResponse)
		if err != nil {
			t.Fatalf("Failed to serialize expected response\n%v", err)
		}

		request, _ := http.NewRequest(http.MethodGet, "/portfolios", nil)
		response := httptest.NewRecorder()

		server := getTestServer(serializedExpectedResponse)

		server.ServeHTTP(response, request)

		var actualResponse *types.PortfolioResponse
		json.Unmarshal(response.Body.Bytes(), &actualResponse)

		if actualResponse == nil {
			t.Fatalf("Response is nil")
		}

		if len(actualResponse.Portfolios) != len(expectedResponse.Portfolios) {
			t.Errorf("Expected response with %d portfolios, got response with %d portfolios", len(expectedResponse.Portfolios), len(actualResponse.Portfolios))
		}

		actualPortfolio1 := actualResponse.Portfolios[0]
		actualPortfolio2 := actualResponse.Portfolios[1]

		assertStringEquals(portfolio1.Name, actualPortfolio1.Name, t)
		assertStringEquals(portfolio1.Type, actualPortfolio1.Type, t)
		assertStringEquals(portfolio1.Uuid, actualPortfolio1.Uuid, t)
		assertStringEquals(portfolio2.Name, actualPortfolio2.Name, t)
		assertStringEquals(portfolio2.Type, actualPortfolio2.Type, t)
		assertStringEquals(portfolio2.Uuid, actualPortfolio2.Uuid, t)

		if portfolio1.Deleted != actualPortfolio1.Deleted {
			t.Errorf("Expected deleted: %t actual deleted: %t", portfolio1.Deleted, actualPortfolio1.Deleted)
		}

		if portfolio2.Deleted != actualPortfolio2.Deleted {
			t.Errorf("Expected deleted: %t actual deleted: %t", portfolio2.Deleted, actualPortfolio2.Deleted)
		}
	})
}

func assertStringEquals(expected, actual string, t *testing.T) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected: %q, Actual: %q", expected, actual)
	}
}
