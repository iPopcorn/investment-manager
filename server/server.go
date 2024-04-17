package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/types"
)

type InvestmentManagerHTTPServer struct {
	client infrastructure.InvestmentManagerExternalHttpClient
	state  *types.State
}

func (s *InvestmentManagerHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %v\n", r)

	route, args := getRouteAndArgsFromPath(r.URL.Path)

	switch route {
	case string(types.Portfolios):
		s.handlePortfolio(w, r, args)
		return
	case string(types.ExecuteStrategy):
		s.handleExecuteStrategy(w, r, args)
		return
	default:
		log.Printf("Route not found: %q\n", route)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func GetDefaultInvestmentManagerHTTPServer() *InvestmentManagerHTTPServer {
	httpClient := infrastructure.GetInvestmentManagerExternalHttpClient()

	return &InvestmentManagerHTTPServer{
		client: *httpClient,
		state: &types.State{
			LastUpdated: time.Now().Format(time.RFC3339),
			Portfolios:  nil,
		},
	}
}

func InvestmentManagerHttpServerFactory(httpClient *infrastructure.InvestmentManagerExternalHttpClient, initialState *types.State) *InvestmentManagerHTTPServer {
	return &InvestmentManagerHTTPServer{
		client: *httpClient,
		state:  initialState,
	}
}

func (s *InvestmentManagerHTTPServer) handlePortfolio(w http.ResponseWriter, r *http.Request, args []string) {
	w.Header().Set("Content-Type", "application/json")
	url := "https://api.coinbase.com/api/v3/brokerage/portfolios"

	if r.Method == http.MethodPost {
		body := r.Body

		defer body.Close()

		bodyData, err := ioutil.ReadAll(body)

		if err != nil {
			log.Printf("Failed to read body from request: %v\n", err)
			writeResponse(w, nil, err)
		}

		resp, err := s.client.Post(url, bodyData)

		writeResponse(w, resp, err)
	} else {
		if len(args) == 1 {
			portfolioUUID := args[0]
			url = url + "/" + portfolioUUID
			resp, err := s.client.Get(url)

			if err != nil {
				log.Printf("Error retrieving portfolio details from URL: %q\nError: %v", url, err)
			}

			writeResponse(w, resp, err)
		} else {
			resp, err := s.client.Get(url)

			if err != nil {
				log.Printf("Error retrieving portfolios from URL: %q\nError: %v", url, err)
			}

			writeResponse(w, resp, err)
		}
	}

	log.Printf("Request handled successfully!")
}

func (s *InvestmentManagerHTTPServer) handleExecuteStrategy(w http.ResponseWriter, r *http.Request, args []string) {
	if r.Method != http.MethodPost {
		err := fmt.Errorf("Invalid http method, wanted %s got %s", http.MethodPost, r.Method)

		writeResponse(w, nil, err)
		return
	}

	now := time.Now()
	// TODO: Use correct last updated
	s.state.LastUpdated = now.Add(time.Second * 5).Format(time.RFC3339)

	body := r.Body

	defer body.Close()

	bodyData, err := ioutil.ReadAll(body)

	if err != nil {
		log.Printf("Failed to read body from request: %v\n", err)
		writeResponse(w, nil, err)
	}

	var requestBody types.ExecuteStrategyRequest

	err = json.Unmarshal(bodyData, &requestBody)

	s.state.Portfolios = []types.Portfolio{
		{
			Name:    requestBody.Portfolio,
			Uuid:    "test-portfolio-id", // TODO: Get portfolio data from coinbase
			Type:    "test",              // TODO: Get portfolio data from coinbase
			Deleted: false,               // TODO: Get portfolio data from coinbase
			CurrentStrategy: &types.Strategy{
				Name:     requestBody.Strategy,
				Currency: requestBody.Currency,
				OpenOffers: []types.Offer{
					{
						ClientOrderId: "test",    // TODO: Generate client order id
						ProductId:     "GBP-ETH", // TODO: Confirm product id + get from request
						Side:          types.BUY,
						Config: types.OrderConfiguration{
							Type:       types.LimitLimitGTD,
							BaseSize:   "10",                                          // TODO: Use correct base size
							LimitPrice: "10",                                          // TODO: Use correct limit price
							PostOnly:   true,                                          // TODO: set post only conditionally
							EndTime:    now.Add(time.Minute * 5).Format(time.RFC3339), // TODO: set conditionally
						},
						SelfTradePreventionId: "test",              // TODO: use correct value
						RetailPortfolioId:     "test-portfolio-id", // TODO confirm correct + use value from actual portfolio
					},
				},
				ClosedOffers: nil,
			},
			PreviousStrategies: nil,
		},
	}

	writeResponse(w, []byte("OK"), nil)
}

func writeResponse(w http.ResponseWriter, response []byte, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Response: \n%s", string(response))

	_, err = w.Write(response)

	if err != nil {
		log.Println("Failed to write response to writer")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getRouteAndArgsFromPath(path string) (string, []string) {
	rawPath := strings.TrimPrefix(path, "/")
	log.Printf("rawPath: %v", rawPath)
	pathTokens := strings.Split(rawPath, "/")
	log.Printf("pathTokens: %v", pathTokens)
	route := pathTokens[0]
	log.Printf("Route: %q\n", route)
	args := []string{}

	for i := 1; i < len(pathTokens); i++ {
		args = append(args, pathTokens[i])
	}

	return route, args
}
