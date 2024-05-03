package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fossoreslp/go-uuid-v4"
	"github.com/iPopcorn/investment-manager/infrastructure"
	"github.com/iPopcorn/investment-manager/server/server_utils"
	"github.com/iPopcorn/investment-manager/server/state"
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
		server_utils.WriteResponse(args.Writer, nil, fmt.Errorf(handlerName+"Invalid http method, wanted %s got %s", http.MethodPost, args.Req.Method))

		return
	}

	body := args.Req.Body

	defer body.Close()

	bodyData, err := ioutil.ReadAll(body)

	if err != nil {
		log.Printf(handlerName+"Failed to read body from request: %v\n", err)
		server_utils.WriteResponse(args.Writer, nil, err)

		return
	}

	var requestBody types.ExecuteStrategyRequest

	err = json.Unmarshal(bodyData, &requestBody)

	if err != nil {
		log.Printf(handlerName + "Failed to deserialize request")
		server_utils.WriteResponse(args.Writer, nil, err)

		return
	}

	userPortfolios, err := server_utils.ListPortfolios(args.Client)

	if err != nil || len(userPortfolios.Portfolios) == 0 {
		log.Printf(handlerName + "Failed to get portfolios from coinbase")
		server_utils.WriteResponse(args.Writer, nil, err)

		return
	}

	var selectedPortfolio *types.Portfolio
	for _, p := range userPortfolios.Portfolios {
		if p.Name == requestBody.Portfolio {
			selectedPortfolio = &p
			break
		}
	}

	if selectedPortfolio == nil {
		server_utils.WriteResponse(args.Writer, nil, fmt.Errorf(handlerName+"Could not find requested portfolio\nGiven: %q", requestBody.Portfolio))

		return
	}

	selectedPortfolioDetails, err := server_utils.PortfolioDetails(args.Client, selectedPortfolio.Uuid)

	if err != nil {
		fmt.Printf("err: %+v\n", err)
		server_utils.WriteResponse(args.Writer, nil, fmt.Errorf(handlerName+"Failed to get details for selected portfolio\n"))

		return
	}

	// TODO: This channel is used for tests, there's probably a better way to handle this
	var finished chan bool
	if len(args.Channels) == 1 {
		finished = args.Channels[0]
	} else {
		finished = make(chan bool)
	}

	productID, err := server_utils.GetProductID(args.Client, selectedPortfolioDetails, string(requestBody.Currency))

	if err != nil {
		fmt.Printf("Error: %+v", err)
		server_utils.WriteResponse(args.Writer, nil, fmt.Errorf(handlerName+"Failed to get ProductID\n"))

		return
	}

	executeStrategyArgs := executeStrategyArgs{
		Client:           args.Client,
		PortfolioDetails: selectedPortfolioDetails,
		StateRepository:  args.StateRepository,
		ProductID:        productID,
		StrategyName:     requestBody.Strategy,
		StrategyCurrency: requestBody.Currency,
		Finished:         finished,
	}

	go executeStrategy(executeStrategyArgs)

	server_utils.WriteResponse(args.Writer, []byte("OK"), nil)
}

type executeStrategyArgs struct {
	Client           *infrastructure.InvestmentManagerExternalHttpClient
	PortfolioDetails *types.PortfolioDetailsResponse
	StateRepository  *state.StateRepository
	ProductID        string
	StrategyName     types.StrategyName
	StrategyCurrency types.SupportedCurrency
	Finished         chan bool
}

func executeStrategy(args executeStrategyArgs) {
	fmt.Println("BEGIN executeStrategy()")
	var newState *types.State
	newState, err := args.StateRepository.GetState()

	if err != nil {
		if strings.Contains(fmt.Sprintf("%s", err), "no such file or directory") {
			fmt.Printf("No state found, initializing\n")

			newState = args.StateRepository.InitState()
			err = args.StateRepository.Save(*newState)

			if err != nil {
				fmt.Printf("failed to init state, returning\n%v\n", err)
				args.Finished <- true
				return
			}
		} else {
			fmt.Printf("Failed to get state from repository\n%v\nReturning\n", err)
			args.Finished <- true
			return
		}
	}

	if err != nil {
		fmt.Printf("Failed to generate uuid for clientOrderId\n%v\nReturning\n", err)

		args.Finished <- true
		return
	}

	breakdown := args.PortfolioDetails.Breakdown
	portfolio := breakdown.Portfolio

	bestBidAsk, err := server_utils.GetBestBidAsk(args.Client, args.ProductID)

	if err != nil {
		fmt.Printf("Failed to get best bid/ask \n%v\n", err)

		args.Finished <- true
		return
	}

	fmt.Printf("Best bid/ask for %s\n%+v\n", args.ProductID, bestBidAsk)

	orderConfigArgs := &createOrderConfigArgs{
		Breakdown:    &breakdown,
		StrategyName: args.StrategyName,
		BestBidAsk:   bestBidAsk,
	}

	orderConfig, err := createOrderConfig(orderConfigArgs)

	if err != nil {
		fmt.Printf("Failed to get order config\n%v\n", err)

		args.Finished <- true
		return
	}

	clientOrderID, err := uuid.NewString()
	newOffer := &types.Offer{
		ClientOrderId:         clientOrderID,
		ProductId:             args.ProductID,
		Side:                  types.BUY,
		Config:                *orderConfig,
		SelfTradePreventionId: types.Default,
		RetailPortfolioId:     portfolio.Uuid,
	}

	previewMode := false

	_, err = server_utils.PlaceOrder(&server_utils.PlaceOrderArgs{
		Client:  args.Client,
		Offer:   newOffer,
		Preview: previewMode,
	})

	if err != nil {
		fmt.Printf("Failed to place order\n%v\n", err)

		args.Finished <- true
		return
	}

	// TODO: Amend state rather than overwrite it
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
					*newOffer,
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

type createOrderConfigArgs struct {
	Breakdown    *types.Breakdown
	StrategyName types.StrategyName
	BestBidAsk   *types.BestBidAskResponse
}

func createOrderConfig(args *createOrderConfigArgs) (*types.OrderConfiguration, error) {
	fiveMinutesFromNow := time.Now().Add(time.Minute * 5).Format(time.RFC3339)

	fmt.Printf("%+v\n", args.BestBidAsk)
	// get available funds, assume GBP for now.
	var availableToTrade float64
	for _, position := range args.Breakdown.SpotPositions {
		if position.Asset == "GBP" {
			availableToTrade = position.AvailableToTradeFiat
			break
		}
	}

	//subtract expected commission
	commissionRate := 0.00400001 // maker commission is 0.40% for orders less than $10k add 0.000001% padding
	expectedCommission := availableToTrade * commissionRate
	availableToTrade -= expectedCommission

	// Match best bid price
	limitPrice := args.BestBidAsk.PriceBooks[0].Bids[0].Price
	limitPriceFloat, err := strconv.ParseFloat(limitPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert limit price to float\nGiven: %q\n%v\n", limitPrice, err)
	}

	// Base size is the quantity of the base currency to buy.
	// Base currency is on the left side of the product id.
	// Example: "ETH-GBP" the base currency is "ETH"
	baseSizeFloat := availableToTrade / limitPriceFloat
	baseSize := strconv.FormatFloat(baseSizeFloat, 'f', 8, 64)

	switch args.StrategyName {
	case types.HODL:
		// For now, HODL strategy will spend all the fiat currency in 1 order
		return &types.OrderConfiguration{
			LimitLimitGTD: types.LimitLimitGTD{
				BaseSize:   baseSize,
				LimitPrice: limitPrice,
				PostOnly:   true,               // TODO: set conditionally
				EndTime:    fiveMinutesFromNow, // TODO: set conditionally
			},
		}, nil
	default:
		return nil, fmt.Errorf("Unsupported strategy name: %q\n", string(args.StrategyName))
	}
}
