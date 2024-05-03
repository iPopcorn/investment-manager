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
	"github.com/iPopcorn/investment-manager/server/server_utils"
	"github.com/iPopcorn/investment-manager/server/state"
	"github.com/iPopcorn/investment-manager/types"
)

type testHttpClient struct {
	getResponseMap map[string][]byte
}

type testReader struct {
	data []byte
}

type testServerArgs struct {
	expectedResponseMap map[string][]byte
	mockRepo            *state.StateRepository
	chans               []chan bool
}

func (testClient testHttpClient) Do(req *http.Request) (*http.Response, error) {
	_, args := server_utils.GetRouteAndArgsFromPath(req.URL.Path)

	index := len(args) - 1
	key := args[index]
	responseData := testClient.getResponseMap[key]

	resp := http.Response{
		Body: io.NopCloser(bytes.NewReader(responseData)),
	}

	return &resp, nil
}

func getTestServer(args *testServerArgs) *InvestmentManagerHTTPServer {
	testHttpClient := testHttpClient{
		getResponseMap: args.expectedResponseMap,
	}
	testInvestmentManagerHTTPClient := infrastructure.InvestmentManagerExternalHttpClient{
		HttpClient: testHttpClient,
	}

	serverArgs := InvestmentManagerHTTPServerArgs{
		HttpClient:      &testInvestmentManagerHTTPClient,
		StateRepository: args.mockRepo,
		Channels:        args.chans,
	}

	return InvestmentManagerHttpServerFactory(serverArgs)
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

		responseMap := make(map[string][]byte)
		responseMap["portfolios"] = serializedExpectedResponse

		testServerArgs := &testServerArgs{
			expectedResponseMap: responseMap,
		}

		server := getTestServer(testServerArgs)

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
		testServerArgs := &testServerArgs{}

		server := getTestServer(testServerArgs)

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

		testZeroBalance := types.Balance{
			Value:    "0",
			Currency: "GBP",
		}

		testPortfolioDetailsResponse := &types.PortfolioDetailsResponse{
			Breakdown: types.Breakdown{
				Portfolio: *testPortfolio,
				PortfolioBalances: types.PortfolioBalances{
					TotalBalance:               testZeroBalance,
					TotalFuturesBalance:        testZeroBalance,
					TotalCashEquivalentBalance: testZeroBalance,
					TotalCryptoBalance:         testZeroBalance,
					FuturesUnrealizedPnl:       testZeroBalance,
					PerpUnrealizedPnl:          testZeroBalance,
				},
				SpotPositions: []types.SpotPositions{
					{
						Asset:                "GBP",
						AccountUuid:          "test-account-uuid",
						TotalBalanceFiat:     0,
						TotalBalanceCrypto:   0,
						AvailableToTradeFiat: 0,
						Allocation:           0,
						CostBasis:            testZeroBalance,
						AssetImgUrl:          "",
						IsCash:               true,
					},
				},
			},
		}

		testProductResponse := &types.ProductResponse{
			Products: []types.Product{
				{
					ProductID: "ETH-GBP",
					Price:     "10",
				},
			},
		}

		testBestBidAskResponse := &types.BestBidAskResponse{
			PriceBooks: []types.PriceBook{
				{
					ProductID: "ETH-GBP",
					Bids: []types.Bid{
						{
							Price: "2349.55",
							Size:  "0.0675",
						},
					},
					Asks: []types.Bid{
						{
							Price: "2350.99",
							Size:  "0.05",
						},
					},
					Time: "2024-05-01T20:07:23.044653Z",
				},
			},
		}

		serializedPortfolioResponse, err := json.Marshal(testPortfolioResponse)
		if err != nil {
			t.Fatalf("Failed to create mock response for test server\n%v", err)
		}

		serializedPortfolioDetailsResponse, err := json.Marshal(testPortfolioDetailsResponse)

		if err != nil {
			t.Fatalf("Failed to create mock response for test server\n%v", err)
		}

		serializedTestProductResponse, err := json.Marshal(testProductResponse)

		if err != nil {
			t.Fatalf("Failed to create mock test product response for test server\n%v", err)
		}

		serializedBestBidAskResponse, err := json.Marshal(testBestBidAskResponse)

		tradeSuccessResponse := make(map[string]string)
		tradeSuccessResponse["order_total"] = "100"
		serializedTestSuccessResponse, err := json.Marshal(tradeSuccessResponse)

		if err != nil {
			t.Fatalf("Failed to serialize data. Given: %+v\n%v\n", tradeSuccessResponse, err)
		}

		timeStart := time.Now()
		formattedTimeStart := timeStart.Format(time.RFC3339)

		testStateRepo := state.StateRepositoryFactory("test-state.json")

		strategyExecutedChannel := make(chan bool)
		testChans := []chan bool{strategyExecutedChannel}

		responseMap := make(map[string][]byte)
		responseMap["portfolios"] = serializedPortfolioResponse
		responseMap["test-portfolio-id"] = serializedPortfolioDetailsResponse
		responseMap["products"] = serializedTestProductResponse
		responseMap["best_bid_ask"] = serializedBestBidAskResponse
		responseMap["preview"] = serializedTestSuccessResponse
		responseMap["orders"] = serializedTestSuccessResponse

		testServerArgs := &testServerArgs{
			expectedResponseMap: responseMap,
			mockRepo:            testStateRepo,
			chans:               testChans,
		}

		testServer := getTestServer(testServerArgs)

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
		if response.Code == http.StatusInternalServerError {
			t.Fatalf("Server did not return 200\n")
		}

		// wait for async operations
		<-strategyExecutedChannel

		updatedState, err := testStateRepo.GetState()
		if err != nil {
			t.Fatalf("Failed to retrieve state from repository")
		}

		if updatedState.LastUpdated == formattedTimeStart {
			t.Errorf(unexpectedUpdate + "LastUpdated equals start time")
		}

		if updatedState.Portfolios == nil {
			t.Fatalf(unexpectedUpdate + "Portfolios is nil")
		}

		if len(updatedState.Portfolios) == 0 {
			t.Fatalf(unexpectedUpdate + "No portfolios added to state")
		}

		if len(updatedState.Portfolios) > 1 {
			t.Errorf(unexpectedUpdate+"expected 1 portfolio but found %d", len(updatedState.Portfolios))
		}

		actualPortfolio := updatedState.Portfolios[0]
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

		assertStringEquals("HODL", string(actualCurrentStrategy.Name), t)
		assertStringEquals("ETH", string(actualCurrentStrategy.Currency), t)

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

		assertStringEquals("ETH-GBP", actualOpenOffer.ProductId, t)
		assertStringEquals(string(types.BUY), string(actualOpenOffer.Side), t)

		offerConfig := actualOpenOffer.Config

		if offerConfig.LimitLimitGTD.BaseSize == "" {
			t.Errorf(unexpectedUpdate + "base size is empty")
		}

		if offerConfig.LimitLimitGTD.LimitPrice == "" {
			t.Errorf(unexpectedUpdate + "limit price is empty")
		}

		if !offerConfig.LimitLimitGTD.PostOnly {
			t.Errorf(unexpectedUpdate + "post only should be true")
		}

		expectedEndTime := timeStart.Add(time.Minute * 5).Format(time.RFC3339)
		assertStringEquals(expectedEndTime, offerConfig.LimitLimitGTD.EndTime, t)

		// TODO: Match the current best bid
		// TODO: Add the transaction to the state
		// TODO: wait for the offer to be fulfilled
		// TODO: update the state
	})
}

func TestTransferFunds(t *testing.T) {
	setup := func(senderAvailableToTradeFiat float64, t *testing.T) (senderID, receiverID string, testServer *InvestmentManagerHTTPServer) {
		t.Helper()
		senderPortfolioID := "test-sender-portfolio-id"
		receiverPortfolioID := "test-receiver-portfolio-id"
		testSenderPortfolio := &types.Portfolio{
			Name:               "test-sender",
			Uuid:               senderPortfolioID,
			Type:               "test",
			Deleted:            false,
			CurrentStrategy:    nil,
			PreviousStrategies: nil,
		}
		testReceiverPortfolio := &types.Portfolio{
			Name:               "test-receiver",
			Uuid:               receiverPortfolioID,
			Type:               "test",
			Deleted:            false,
			CurrentStrategy:    nil,
			PreviousStrategies: nil,
		}
		testZeroBalance := types.Balance{
			Value:    "0",
			Currency: "GBP",
		}

		testSenderPortfolioDetailsResponse := &types.PortfolioDetailsResponse{
			Breakdown: types.Breakdown{
				Portfolio: *testSenderPortfolio,
				PortfolioBalances: types.PortfolioBalances{
					TotalBalance:               testZeroBalance,
					TotalFuturesBalance:        testZeroBalance,
					TotalCashEquivalentBalance: testZeroBalance,
					TotalCryptoBalance:         testZeroBalance,
					FuturesUnrealizedPnl:       testZeroBalance,
					PerpUnrealizedPnl:          testZeroBalance,
				},
				SpotPositions: []types.SpotPositions{
					{
						Asset:                "GBP",
						AccountUuid:          "test-account-uuid",
						TotalBalanceFiat:     0,
						TotalBalanceCrypto:   0,
						AvailableToTradeFiat: senderAvailableToTradeFiat,
						Allocation:           0,
						CostBasis:            testZeroBalance,
						AssetImgUrl:          "",
						IsCash:               true,
					},
				},
			},
		}

		testReceiverPortfolioDetailsResponse := &types.PortfolioDetailsResponse{
			Breakdown: types.Breakdown{
				Portfolio: *testReceiverPortfolio,
				PortfolioBalances: types.PortfolioBalances{
					TotalBalance:               testZeroBalance,
					TotalFuturesBalance:        testZeroBalance,
					TotalCashEquivalentBalance: testZeroBalance,
					TotalCryptoBalance:         testZeroBalance,
					FuturesUnrealizedPnl:       testZeroBalance,
					PerpUnrealizedPnl:          testZeroBalance,
				},
				SpotPositions: []types.SpotPositions{
					{
						Asset:                "GBP",
						AccountUuid:          "test-account-uuid",
						TotalBalanceFiat:     0,
						TotalBalanceCrypto:   0,
						AvailableToTradeFiat: 0,
						Allocation:           0,
						CostBasis:            testZeroBalance,
						AssetImgUrl:          "",
						IsCash:               true,
					},
				},
			},
		}

		transferFundsSuccessResponse := make(map[string]string)
		transferFundsSuccessResponse["source_portfolio_uuid"] = senderPortfolioID
		transferFundsSuccessResponse["target_portfolio_uuid"] = receiverPortfolioID

		serializedTransferFundsSuccessResponse, err := json.Marshal(transferFundsSuccessResponse)
		if err != nil {
			t.Fatalf("Failed to serialize data\nGiven: %+v\n%v\n", transferFundsSuccessResponse, err)
		}

		serializedTestSenderPortfolioDetailsResponse, err := json.Marshal(testSenderPortfolioDetailsResponse)
		if err != nil {
			t.Fatalf("Failed to serialize data\nGiven: %+v\n%v\n", testSenderPortfolioDetailsResponse, err)
		}

		serializedTestReceiverPortfolioDetailsResponse, err := json.Marshal(testReceiverPortfolioDetailsResponse)
		if err != nil {
			t.Fatalf("Failed to serialize data\nGiven: %+v\n%v\n", testReceiverPortfolioDetailsResponse, err)
		}

		responseMap := make(map[string][]byte)
		responseMap[senderPortfolioID] = serializedTestSenderPortfolioDetailsResponse
		responseMap[receiverPortfolioID] = serializedTestReceiverPortfolioDetailsResponse
		responseMap["move_funds"] = serializedTransferFundsSuccessResponse

		testServerArgs := &testServerArgs{
			expectedResponseMap: responseMap,
			mockRepo:            nil,
			chans:               nil,
		}

		return senderPortfolioID, receiverPortfolioID, getTestServer(testServerArgs)
	}
	t.Run("Fails to transfer if not enough funds", func(t *testing.T) {
		// Arrange
		senderID, receiverID, testServer := setup(0, t)
		body := types.TransferRequest{
			SenderID:   senderID,
			ReceiverID: receiverID,
			Amount:     "10",
		}

		serializedBody, err := json.Marshal(body)

		if err != nil {
			t.Fatalf("Failed to create body for request\n%v", err)
		}

		request, err := http.NewRequest(http.MethodPost, "/"+string(types.TransferFunds), bytes.NewReader(serializedBody))

		if err != nil {
			t.Fatalf("Failed to create http request\n%v", err)
		}

		response := httptest.NewRecorder()

		// Act
		testServer.ServeHTTP(response, request)

		// Assert
		if response.Code != http.StatusBadRequest {
			t.Fatalf("Expected %d Received %d", http.StatusBadRequest, response.Code)
		}
	})

	t.Run("Transfers funds as expected", func(t *testing.T) {
		// Arrange
		senderID, receiverID, testServer := setup(20, t)
		body := types.TransferRequest{
			SenderID:   senderID,
			ReceiverID: receiverID,
			Amount:     "10",
		}

		serializedBody, err := json.Marshal(body)

		if err != nil {
			t.Fatalf("Failed to create body for request\n%v", err)
		}

		request, err := http.NewRequest(http.MethodPost, "/"+string(types.TransferFunds), bytes.NewReader(serializedBody))

		if err != nil {
			t.Fatalf("Failed to create http request\n%v", err)
		}

		response := httptest.NewRecorder()

		// Act
		testServer.ServeHTTP(response, request)

		// Assert
		if response.Code != http.StatusOK {
			t.Fatalf("Expected %d Received %d", http.StatusOK, response.Code)
		}
	})
}

func assertStringEquals(expected, actual string, t *testing.T) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected: %q, Actual: %q", expected, actual)
	}
}
