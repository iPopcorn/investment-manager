package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/server/state"
	"github.com/iPopcorn/investment-manager/server/util"
	"github.com/iPopcorn/investment-manager/types"
)

type HandleExecuteStrategyArgs struct {
	Client          *infrastructure.InvestmentManagerExternalHttpClient
	Writer          http.ResponseWriter
	Req             *http.Request
	Args            []string
	Channels        []chan bool
	StateRepository *state.StateRepository
}

func HandleExecuteStrategy(args HandleExecuteStrategyArgs) {
	handlerName := "handleExecuteStrategy: "

	if args.Req.Method != http.MethodPost {
		util.WriteResponse(args.Writer, nil, fmt.Errorf(handlerName+"Invalid http method, wanted %s got %s", http.MethodPost, args.Req.Method))

		return
	}

	body := args.Req.Body

	defer body.Close()

	bodyData, err := ioutil.ReadAll(body)

	if err != nil {
		log.Printf(handlerName+"Failed to read body from request: %v\n", err)
		util.WriteResponse(args.Writer, nil, err)

		return
	}

	var requestBody types.ExecuteStrategyRequest

	err = json.Unmarshal(bodyData, &requestBody)

	if err != nil {
		log.Printf(handlerName + "Failed to deserialize request")
		util.WriteResponse(args.Writer, nil, err)

		return
	}

	userPortfolios, err := util.ListPortfolios(args.Client)

	if err != nil || len(userPortfolios.Portfolios) == 0 {
		log.Printf(handlerName + "Failed to get portfolios from coinbase")
		util.WriteResponse(args.Writer, nil, err)

		return
	}

	var selectedPortfolio *types.Portfolio

	for _, p := range userPortfolios.Portfolios {
		if p.Name == requestBody.Portfolio {
			selectedPortfolio = &p
		}
	}

	if selectedPortfolio == nil {
		util.WriteResponse(args.Writer, nil, fmt.Errorf(handlerName+"Could not find requested portfolio\nGiven: %q", requestBody.Portfolio))

		return
	}

	// TODO: This channel is used for tests, there's probably a better way to handle this
	var finished chan bool
	if len(args.Channels) == 1 {
		finished = args.Channels[0]
	} else {
		finished = make(chan bool)
	}

	go executeStrategy(*selectedPortfolio, args.StateRepository, requestBody, finished)

	util.WriteResponse(args.Writer, []byte("OK"), nil)
}

func executeStrategy(portfolio types.Portfolio, stateRepository *state.StateRepository, executeStrategyRequest types.ExecuteStrategyRequest, finished chan bool) {
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

	stateRepository.Save(*newState)
	fmt.Printf("END executeStrategy()\n")
	finished <- true
}
