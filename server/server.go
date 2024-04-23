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
	"github.com/iPopcorn/investment-manager/server/util"
	"github.com/iPopcorn/investment-manager/types"
)

type InvestmentManagerHTTPServer struct {
	client          infrastructure.InvestmentManagerExternalHttpClient
	stateRepository *StateRepository
	channels        []chan bool
}

type InvestmentManagerHTTPServerArgs struct {
	HttpClient      *infrastructure.InvestmentManagerExternalHttpClient
	StateRepository *StateRepository
	Channels        []chan bool
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

	initialState := &types.State{
		LastUpdated: time.Now().Format(time.RFC3339),
		Portfolios:  nil,
	}

	stateRepo := StateRepositoryFactory(*initialState)

	return &InvestmentManagerHTTPServer{
		client:          *httpClient,
		stateRepository: stateRepo,
	}
}

func InvestmentManagerHttpServerFactory(args InvestmentManagerHTTPServerArgs) *InvestmentManagerHTTPServer {
	return &InvestmentManagerHTTPServer{
		client:          *args.HttpClient,
		stateRepository: args.StateRepository,
		channels:        args.Channels,
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
			util.WriteResponse(w, nil, err)
		}

		resp, err := s.client.Post(url, bodyData)

		util.WriteResponse(w, resp, err)
	} else {
		if len(args) == 1 {
			portfolioUUID := args[0]
			url = url + "/" + portfolioUUID
			resp, err := s.client.Get(url)

			if err != nil {
				log.Printf("Error retrieving portfolio details from URL: %q\nError: %v", url, err)
			}

			util.WriteResponse(w, resp, err)
		} else {
			resp, err := s.client.Get(url)

			if err != nil {
				log.Printf("Error retrieving portfolios from URL: %q\nError: %v", url, err)
			}

			util.WriteResponse(w, resp, err)
		}
	}

	log.Printf("Request handled successfully!")
}

func (s *InvestmentManagerHTTPServer) handleExecuteStrategy(w http.ResponseWriter, r *http.Request, args []string) {
	handlerName := "handleExecuteStrategy: "

	if r.Method != http.MethodPost {
		util.WriteResponse(w, nil, fmt.Errorf(handlerName+"Invalid http method, wanted %s got %s", http.MethodPost, r.Method))

		return
	}

	body := r.Body

	defer body.Close()

	bodyData, err := ioutil.ReadAll(body)

	if err != nil {
		log.Printf(handlerName+"Failed to read body from request: %v\n", err)
		util.WriteResponse(w, nil, err)

		return
	}

	var requestBody types.ExecuteStrategyRequest

	err = json.Unmarshal(bodyData, &requestBody)

	if err != nil {
		log.Printf(handlerName + "Failed to deserialize request")
		util.WriteResponse(w, nil, err)

		return
	}

	userPortfolios, err := util.ListPortfolios(&s.client)

	if err != nil || len(userPortfolios.Portfolios) == 0 {
		log.Printf(handlerName + "Failed to get portfolios from coinbase")
		util.WriteResponse(w, nil, err)

		return
	}

	var selectedPortfolio *types.Portfolio

	for _, p := range userPortfolios.Portfolios {
		if p.Name == requestBody.Portfolio {
			selectedPortfolio = &p
		}
	}

	if selectedPortfolio == nil {
		util.WriteResponse(w, nil, fmt.Errorf(handlerName+"Could not find requested portfolio\nGiven: %q", requestBody.Portfolio))

		return
	}

	var finished chan bool
	if len(s.channels) == 1 {
		finished = s.channels[0]
	} else {
		finished = make(chan bool)
	}

	go executeStrategy(*selectedPortfolio, s.stateRepository, requestBody, finished)

	util.WriteResponse(w, []byte("OK"), nil)
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

// Not sure if state repo should be ref or copy - probably ref until state is persisted somehow
func executeStrategy(portfolio types.Portfolio, stateRepository *StateRepository, executeStrategyRequest types.ExecuteStrategyRequest, finished chan bool) {
	fmt.Println("BEGIN executeStrategy()")
	newState, err := stateRepository.GetState()

	if err != nil {
		fmt.Printf("Failed to get state from repository\n%v\nReturning\n", err)
		return
	}

	fiveMinutesFromNow := time.Now().Add(time.Minute * 5).Format(time.RFC3339)

	newState.Portfolios = []types.Portfolio{
		{
			Name:    portfolio.Name,
			Uuid:    portfolio.Uuid,
			Type:    portfolio.Type,
			Deleted: portfolio.Deleted,
			CurrentStrategy: &types.Strategy{
				Name:     executeStrategyRequest.Strategy,
				Currency: executeStrategyRequest.Currency,
				OpenOffers: []types.Offer{
					{
						ClientOrderId: "test",    // TODO: Generate client order id
						ProductId:     "GBP-ETH", // TODO: Confirm product id + get from request
						Side:          types.BUY,
						Config: types.OrderConfiguration{
							Type:       types.LimitLimitGTD,
							BaseSize:   "10",               // TODO: Use correct base size
							LimitPrice: "10",               // TODO: Use correct limit price
							PostOnly:   true,               // TODO: set post only conditionally
							EndTime:    fiveMinutesFromNow, // TODO: set conditionally
						},
						SelfTradePreventionId: "test",         // TODO: use correct value
						RetailPortfolioId:     portfolio.Uuid, // TODO confirm correct
					},
				},
				ClosedOffers: nil,
			},
			PreviousStrategies: nil,
		},
	}

	newState.LastUpdated = time.Now().Add(time.Second).Format(time.RFC3339)

	stateRepository.Save(newState)
	fmt.Printf("END executeStrategy()\n")
	finished <- true
}
