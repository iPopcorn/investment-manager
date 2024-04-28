package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fossoreslp/go-uuid-v4"
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

	selectedPortfolioDetails, err := util.PortfolioDetails(args.Client, selectedPortfolio.Uuid)

	if err != nil {
		fmt.Printf("err: %+v\n", err)
		util.WriteResponse(args.Writer, nil, fmt.Errorf(handlerName+"Failed to get details for selected portfolio\n"))

		return
	}

	// TODO: This channel is used for tests, there's probably a better way to handle this
	var finished chan bool
	if len(args.Channels) == 1 {
		finished = args.Channels[0]
	} else {
		finished = make(chan bool)
	}

	productID, err := util.GetProductID(args.Client, selectedPortfolioDetails, requestBody.Currency)

	if err != nil {
		fmt.Printf("Error: %+v", err)
		util.WriteResponse(args.Writer, nil, fmt.Errorf(handlerName+"Failed to get ProductID\n"))

		return
	}

	executeStrategyArgs := executeStrategyArgs{
		PortfolioDetails: selectedPortfolioDetails,
		StateRepository:  args.StateRepository,
		ProductID:        productID,
		StrategyName:     requestBody.Strategy,
		StrategyCurrency: requestBody.Currency,
		Finished:         finished,
	}

	go executeStrategy(executeStrategyArgs)

	util.WriteResponse(args.Writer, []byte("OK"), nil)
}

type executeStrategyArgs struct {
	PortfolioDetails *types.PortfolioDetailsResponse
	StateRepository  *state.StateRepository
	ProductID        string
	StrategyName     string
	StrategyCurrency string
	Finished         chan bool
}

func executeStrategy(args executeStrategyArgs) {
	fmt.Println("BEGIN executeStrategy()")
	newState, err := args.StateRepository.GetState()

	if err != nil {
		fmt.Printf("Failed to get state from repository\n%v\nReturning\n", err)

		args.Finished <- true
		return
	}

	fiveMinutesFromNow := time.Now().Add(time.Minute * 5).Format(time.RFC3339)
	clientOrderID, err := uuid.NewString()

	if err != nil {
		fmt.Printf("Failed to generate uuid for clientOrderId\n%v\nReturning\n", err)

		args.Finished <- true
		return
	}

	portfolio := args.PortfolioDetails.Breakdown.Portfolio

	newState.Portfolios = []types.Portfolio{
		{
			Name:    portfolio.Name,
			Uuid:    portfolio.Uuid,
			Type:    portfolio.Type,
			Deleted: portfolio.Deleted,
			CurrentStrategy: &types.Strategy{
				Name:     args.StrategyName,
				Currency: args.StrategyCurrency,
				OpenOffers: []types.Offer{
					{
						ClientOrderId: clientOrderID,
						ProductId:     args.ProductID,
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

	args.StateRepository.Save(*newState)
	fmt.Printf("END executeStrategy()\n")
	args.Finished <- true
}
