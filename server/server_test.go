package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func getTestServer(expectedResponse []byte, initialState *types.State) *InvestmentManagerHTTPServer {
	testHttpClient := testHttpClient{
		getResponse: expectedResponse,
	}
	testInvestmentManagerHTTPClient := infrastructure.InvestmentManagerExternalHttpClient{
		HttpClient: testHttpClient,
	}
	return InvestmentManagerHttpServerFactory(&testInvestmentManagerHTTPClient, initialState)
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

		server := getTestServer(serializedExpectedResponse, nil)

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

	t.Run("Handles path not found", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/portfolio", nil)
		response := httptest.NewRecorder()

		server := getTestServer(nil, nil)

		server.ServeHTTP(response, request)

		actual := response.Result().StatusCode
		expected := http.StatusNotFound

		if actual != expected {
			t.Fatalf("Expected: %d Actual: %d", expected, actual)
		}
	})
}

func TestExecuteStrategy(t *testing.T) {
	t.Run("Executes the HODL strategy", func(t *testing.T) {
		// Arrange
		/*
			beginning state: {
				last_updated: now
				portfolios: []
			}
		*/
		testPortfolio := &types.Portfolio{
			Name:               "test",
			Uuid:               "test-portfolio-id",
			Type:               "test",
			Deleted:            false,
			CurrentStrategy:    nil,
			PreviousStrategies: nil,
		}

		testPortfolioResponse := &types.PortfolioResponse{
			Portfolios: []types.Portfolio{*testPortfolio},
		}

		serializedResponse, err := json.Marshal(testPortfolioResponse)

		if err != nil {
			t.Fatalf("Failed to create mock response for test server\n%v", err)
		}

		timeStart := time.Now()
		formattedTimeStart := timeStart.Format(time.RFC3339)
		initialState := &types.State{
			LastUpdated: formattedTimeStart,
			Portfolios:  nil,
		}

		testServer := getTestServer(serializedResponse, initialState)

		body := types.ExecuteStrategyRequest{
			Portfolio: testPortfolio.Name,
			Strategy:  "HODL",
			Currency:  "ETH",
		}

		serializedBody, err := json.Marshal(body)

		if err != nil {
			t.Fatalf("Failed to create body for request\n%v", err)
		}

		request, err := http.NewRequest(http.MethodPost, "/"+string(types.ExecuteStrategy), bytes.NewReader(serializedBody))

		if err != nil {
			t.Fatalf("Failed to create http request\n%v", err)
		}

		response := httptest.NewRecorder()

		// Act
		testServer.ServeHTTP(response, request)

		// Assert
		// Update state to say that the portfolio is executing a strategy
		/*
			{
				last_updated: timestamp
				portfolios: [
					{
						name: string
						current_strategy: strategy{
							name: string
							currency: string
							open_offers: []offer [
								offer{
									client_order_id: string
									product_id: string ('BTC-USD')
									side: string ('BUY' 'SELL')
									order_configuration: {
										limit_limit_gtd (only use good-to-date so that orders are cancelled if not filled by certain time): {
											base_size: string
											limit_price: string
											post_only: boolean (should be true to indicate maker only (lower fees))
											end_time: RFC3339 Timestamp
										}
									}
									self_trade_prevention_id: string (hardcoded const)
									retail_portfolio_id: string
								}
							]
							closed_offers: []offer
						}
						previous_strategies: []strategy
						...
					}
				]
			}
		*/
		const unexpectedUpdate = "State was not updated as expected, "
		if initialState.LastUpdated == formattedTimeStart {
			t.Errorf(unexpectedUpdate + "LastUpdated equals start time")
		}

		if initialState.Portfolios == nil {
			t.Fatalf(unexpectedUpdate + "Portfolios is nil")
		}

		if len(initialState.Portfolios) == 0 {
			t.Fatalf(unexpectedUpdate + "No portfolios added to state")
		}

		if len(initialState.Portfolios) > 1 {
			t.Errorf(unexpectedUpdate+"expected 1 portfolio but found %d", len(initialState.Portfolios))
		}

		actualPortfolio := initialState.Portfolios[0]
		assertStringEquals(testPortfolio.Name, actualPortfolio.Name, t)
		assertStringEquals(testPortfolio.Uuid, actualPortfolio.Uuid, t)
		assertStringEquals(testPortfolio.Type, actualPortfolio.Type, t)
		if actualPortfolio.Deleted {
			t.Errorf(unexpectedUpdate + "portfolio marked as deleted but should not be")
		}

		if actualPortfolio.CurrentStrategy == nil {
			t.Fatalf(unexpectedUpdate + "no current strategy found")
		}

		if actualPortfolio.PreviousStrategies != nil {
			t.Errorf(unexpectedUpdate + "previous strategies found when none expected")
		}

		actualCurrentStrategy := actualPortfolio.CurrentStrategy

		assertStringEquals("HODL", actualCurrentStrategy.Name, t)
		assertStringEquals("ETH", actualCurrentStrategy.Currency, t)

		if actualCurrentStrategy.OpenOffers == nil {
			t.Fatalf(unexpectedUpdate + "no open offers found")
		}

		if actualCurrentStrategy.ClosedOffers != nil {
			t.Errorf(unexpectedUpdate + "found closed offers when none expected")
		}

		// Place a buy order for the given currency
		if len(actualCurrentStrategy.OpenOffers) < 1 {
			t.Fatalf(unexpectedUpdate + "open offers is empty")
		}

		if len(actualCurrentStrategy.OpenOffers) > 1 {
			t.Errorf(unexpectedUpdate+"expected 1 open offer but found %d", len(actualCurrentStrategy.OpenOffers))
		}

		actualOpenOffer := actualCurrentStrategy.OpenOffers[0]

		if actualOpenOffer.ClientOrderId == "" {
			t.Errorf(unexpectedUpdate + "client order id is empty")
		}

		assertStringEquals("GBP-ETH", actualOpenOffer.ProductId, t)
		assertStringEquals(string(types.BUY), string(actualOpenOffer.Side), t)

		offerConfig := actualOpenOffer.Config

		assertStringEquals(string(types.LimitLimitGTD), string(offerConfig.Type), t)

		if offerConfig.BaseSize == "" {
			t.Errorf(unexpectedUpdate + "base size is empty")
		}

		if offerConfig.LimitPrice == "" {
			t.Errorf(unexpectedUpdate + "limit price is empty")
		}

		if !offerConfig.PostOnly {
			t.Errorf(unexpectedUpdate + "post only should be true")
		}

		expectedEndTime := timeStart.Add(time.Minute * 5).Format(time.RFC3339)
		assertStringEquals(expectedEndTime, offerConfig.EndTime, t)

		// TODO: Match the current best bid
		// TODO: Add the transaction to the state
		// TODO: wait for the offer to be fulfilled
		// TODO: update the state
	})
}

func assertStringEquals(expected, actual string, t *testing.T) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected: %q, Actual: %q", expected, actual)
	}
}
